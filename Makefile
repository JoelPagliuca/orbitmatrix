NAME := caliban

build: $(NAME)

$(NAME): $(wildcard *.go)
	@echo "+ $@"
	go build

test: build
	@echo "+ $@"
	@./scripts/run-tests.py

run: build
	@echo "+ $@"
	@./$(NAME)

clean:
	@echo "+ $@"
	rm -f $(NAME)
