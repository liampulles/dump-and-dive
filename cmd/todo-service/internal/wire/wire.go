package wire

import (
	"fmt"
	"os"

	goConfig "github.com/liampulles/go-config"

	"github.com/liampulles/dump-and-dive/cmd/todo-service/internal/adapter"
	"github.com/liampulles/dump-and-dive/cmd/todo-service/internal/domain"
	"github.com/liampulles/dump-and-dive/cmd/todo-service/internal/driver/goconfig"
	"github.com/liampulles/dump-and-dive/cmd/todo-service/internal/driver/grpc"
	"github.com/liampulles/dump-and-dive/cmd/todo-service/internal/usecase"
)

// Run is the main entrypoint for todo-service
func Run(source goConfig.Source) int {
	server := wire(source)
	if err := server.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err.Error())
		return 1
	}
	return 0
}

func wire(source goConfig.Source) grpc.Server {
	timeProvider := domain.NewGoTimeProvider()

	domainFactory := domain.NewTODOFactoryImpl(timeProvider)
	stateModifier := &adapter.DummyStateModifier{}

	usecaseCommandService := usecase.NewCommandServiceImpl(domainFactory, stateModifier)

	commandController := adapter.NewCommandControllerImpl(usecaseCommandService)

	cfgProvider := goconfig.NewProvider(source)
	grpcCommandService := grpc.NewCommandServiceImpl(commandController)

	return grpc.NewServerImpl(cfgProvider, grpcCommandService)
}
