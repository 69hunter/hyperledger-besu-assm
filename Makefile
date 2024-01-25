all: 
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/hyperledger-besu-assm main.go
	chmod 655 bin/hyperledger-besu-assm
	zip main.zip bin/hyperledger-besu-assm
	mv main.zip artifacts/main.zip