package main

import (
	app "ethereum-grpc/server/internal/app"
)

var port = 44044

var infura_url = "https://mainnet.infura.io/v3/d46bdacab9a64306a97e95f419caebbd"

func main() {
	application := app.New(port)

	application.GRPCServer.Run()
}
