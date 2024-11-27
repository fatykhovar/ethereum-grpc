package app

import (
	grpcapp "ethereum-grpc/server/internal/app/grpc"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(port int) *App {
	grpcApp := grpcapp.New(port)

	return &App{GRPCServer: grpcApp}
}
