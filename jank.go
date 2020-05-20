package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

// JankedHandler implements http.Handler
// F must be one of
// 	* func(ResponseWriter, *Request)
// 	* func(ResponseWriter, *Request, interface{})
// and must output one of
// 	* nil
// 	* interface{}
// 	* error
// 	* (interface{}, error)
type JankedHandler struct {
	F interface{}
}

// ServeHTTP - first figure out what F is and call it
func (hw JankedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// function itself
	fn := reflect.ValueOf(hw.F)
	// function signature
	sig := fn.Type()

	// check that the first two arguments are of the correct type
	if sig.NumIn() < 2 {
		log.Panicln("inner function has wrong amount of args")
	}
	if !(reflect.TypeOf(w).AssignableTo(sig.In(0)) && reflect.TypeOf(r).AssignableTo(sig.In(1))) {
		log.Panicln("First two args of a function were bad")
	}
	argv := make([]reflect.Value, sig.NumIn())
	argv[0] = reflect.ValueOf(w)
	argv[1] = reflect.ValueOf(r)
	// fill up a struct arg from the request input
	if sig.NumIn() > 2 && sig.In(2).Kind() == reflect.Struct {
		input := sig.In(2)
		arg := reflect.New(input).Elem()
		// iterate over all the elements in the third argument
		for i := 0; i < input.NumField(); i++ {
			if !arg.Field(i).CanSet() {
				continue
			}
			fld := input.Field(i)
			from, ok := fld.Tag.Lookup("from")
			if !ok {
				continue
			}
			switch from {
			case "body":
				val := reflect.New(fld.Type)
				defer r.Body.Close()
				err := json.NewDecoder(r.Body).Decode(val.Interface())
				arg.Field(i).Set(val.Elem())
				if err != nil {
					log.Println("Error: ", err.Error())
				}
			}
		}
		argv[2] = arg
	}
	// call the function
	rets := fn.Call(argv)
	// give the output to the ResponseWriter
	for _, r := range rets {
		if r.Kind() == reflect.Struct || r.Kind() == reflect.Slice {
			payload, _ := json.Marshal(r.Interface())
			w.Write(payload)
			break
		}
		if r.IsNil() {
			continue
		}
		if e := (*error)(nil); r.Type().Implements(reflect.TypeOf(e).Elem()) {
			msg := fmt.Sprintf("%v", r)
			payload, err := json.Marshal(struct{ Error string }{msg})
			if err != nil {
				fmt.Println(err.Error())
			}
			w.Write(payload)
		}
	}
}

// AddSwagger adds a swagger-like endpoint to the mux
func AddSwagger(mux *http.ServeMux, rts []Route) {
	log.Println("API description available on /swagger.txt")
	rsp := "API"
	types := make(map[reflect.Type]bool)
	for _, r := range rts {
		rsp += "\n---\n" + r.Pattern
		rsp += "\nauth: " + strconv.FormatBool(!r.Config.AllowAnonymous)
		rsp += " methods: " + strings.Join(r.Config.AllowedMethods, " ")
		fn := reflect.ValueOf(r.HandlerFunc)
		sig := fn.Type()
		if sig.NumIn() >= 3 {
			rsp += "\ninput:"
			a3 := sig.In(2)
			for i := 0; i < a3.NumField(); i++ {
				fld := a3.Field(i)
				if f, ok := fld.Tag.Lookup("from"); ok {
					rsp += " " + fld.Type.String()[5:] + " (" + f + ") "
					types[fld.Type] = true
				}
			}
		}
	}
	rsp += "\n---\n\nTYPES\n---"
	for t := range types {
		n := reflect.New(t)
		b, _ := json.Marshal(n.Interface())
		rsp += "\n" + t.String()[5:]
		rsp += fmt.Sprintf(" %v", string(b))
	}
	mux.HandleFunc("/swagger.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain")
		w.Write(bytes.NewBufferString(rsp).Bytes())
	})
}
