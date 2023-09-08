.DEFAULT_GOAL := example

example: 
	@echo "Frameplay interview."

test: 
	go test -v

run: 
	go build -o app && ./app