run: bin/hmb
	@PATH="$(PWD)/bin:$(PATH)" heroku local

bin/hmb: main.go
	go build -o bin/hmb main.go

clean:
	rm -rf bin