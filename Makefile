all: 
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/main main.go
	chmod 655 bin/main
	zip main.zip bin/main
	mv main.zip artifacts/main.zip