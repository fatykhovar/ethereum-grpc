package agrpcpp

import (
	ethereumGRPC "ethereum-grpc/server/internal/grpc"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type App struct {
	gRPCServer *grpc.Server
	port       int
}

func New(port int) *App {
	gRPCServer := grpc.NewServer()

	ethereumGRPC.Register(gRPCServer)

	return &App{
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Println("grpc server is running	on :", a.port)

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	log.Println("Stopping gRPC server...")

	a.gRPCServer.GracefulStop()
}
