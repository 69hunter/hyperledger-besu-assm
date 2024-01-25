package main

import (
	"github.com/69hunter/hyperledger-besu-assm/app"
	"github.com/69hunter/hyperledger-besu-assm/core"
	"github.com/69hunter/hyperledger-besu-assm/localstorage"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	aCore := core.NewAdapter()
	aLocalStorage := localstorage.NewAdapter()
	app := app.NewAdapter(aCore, aLocalStorage)
	lambda.Start(app.LambdaHandler)
}
