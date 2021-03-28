package grpc

import (
	"context"
	"fmt"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/liampulles/dump-and-dive/cmd/todo-service/internal/adapter"
	"github.com/liampulles/dump-and-dive/cmd/todo-service/internal/domain"
	"github.com/liampulles/dump-and-dive/cmd/todo-service/internal/driver/grpc/gen"
	"github.com/liampulles/dump-and-dive/cmd/todo-service/internal/usecase"
)

// --------------- Server ---------------

// Server serves gRPC requests
type Server interface {
	Run() error
}

// ServerImpl implements Server
type ServerImpl struct {
	configProvider adapter.ConfigProvider
	commandService gen.CommandServiceServer
}

var _ Server = &ServerImpl{}

// NewServerImpl is a constructor
func NewServerImpl(
	configProvider adapter.ConfigProvider,
	commandService gen.CommandServiceServer,
) *ServerImpl {

	return &ServerImpl{
		configProvider: configProvider,
		commandService: commandService,
	}
}

// Run runs the server
// This is a long lived operation, run in it in a goroutine if
// you want to do something else too
func (s *ServerImpl) Run() error {
	snap, err := s.configProvider.GetAdapterConfig()
	if err != nil {
		return fmt.Errorf("could not get config: %w", err)
	}
	lis, err := createListener(snap.Port())
	if err != nil {
		return fmt.Errorf("could not create listener: %w", err)
	}
	server := s.createServer()

	fmt.Fprintf(os.Stderr, "Running on port %d...\n", snap.Port())
	return server.Serve(lis)
}

func (s *ServerImpl) createServer() *grpc.Server {
	server := grpc.NewServer()
	gen.RegisterCommandServiceServer(server, s.commandService)
	return server
}

func createListener(port int) (net.Listener, error) {
	return net.Listen("tcp", fmt.Sprintf(":%d", port))
}

// --------------- Services ---------------

// CommandServiceImpl implements gen.CommandServiceServer
type CommandServiceImpl struct {
	controller adapter.CommandController

	gen.UnimplementedCommandServiceServer
}

// NewCommandServiceImpl is a constructor
func NewCommandServiceImpl(controller adapter.CommandController) *CommandServiceImpl {
	return &CommandServiceImpl{
		controller: controller,
	}
}

var _ gen.CommandServiceServer = &CommandServiceImpl{}

// Create is a gRPC method for creating TODOs
func (cs *CommandServiceImpl) Create(ctx context.Context, request *gen.CreateRequest) (*gen.EntityID, error) {
	usecaseRequest := mapCreateRequest(request)
	id, err := cs.controller.Create(usecaseRequest)
	if err != nil {
		return nil, convertError(err)
	}
	return mapEntityIDOut(id), nil
}

func mapCreateRequest(in *gen.CreateRequest) *usecase.CreateRequest {
	return &usecase.CreateRequest{
		Name:    in.Name,
		Details: in.Details,
		Due:     in.Due.AsTime(),
	}
}

func mapEntityIDOut(in usecase.EntityID) *gen.EntityID {
	return &gen.EntityID{
		Id: uint32(in),
	}
}

func convertError(in error) error {
	switch errV := in.(type) {
	case *domain.FieldError:
		return status.Error(codes.InvalidArgument, errV.Error())

	}
	return status.Error(codes.Unknown, in.Error())
}
