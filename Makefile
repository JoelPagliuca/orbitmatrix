NAME := caliban

build: $(NAME)

$(NAME): $(wildcard *.go)
	@echo "+ $@"
	go build

clean:
	rm -f $(NAME)
