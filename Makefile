NAME := orbitmatrix

build: $(NAME)

$(NAME): $(wildcard *.go)
	@echo "+ $@"
	go build

test: test-unit test-contract

test-unit:
	@echo "+ $@"
	@go test

test-contract: build
	@echo "+ $@"
	@./scripts/run-tests.py

run: build
	@echo "+ $@"
	@./$(NAME)

clean:
	@echo "+ $@"
	rm -f $(NAME)
