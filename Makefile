NAME := caliban

build: $(NAME)

$(NAME): $(wildcard *.go)
	@echo "+ $@"
	go build

test: build
	@echo "+ $@"
	@./scripts/run-tests.sh

clean:
	@echo "+ $@"
	rm -f $(NAME)
