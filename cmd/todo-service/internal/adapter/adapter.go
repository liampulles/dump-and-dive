package adapter

import (
	"errors"

	"github.com/liampulles/dump-and-dive/cmd/todo-service/internal/domain"
	"github.com/liampulles/dump-and-dive/cmd/todo-service/internal/usecase"
)

// --------------- Config ---------------

// ConfigSnapshot gives a view of adapter config properties at a point in time
type ConfigSnapshot interface {
	Port() int
}

// ConfigProvider provides config
type ConfigProvider interface {
	GetAdapterConfig() (ConfigSnapshot, error)
}

// --------------- Controllers ---------------

// CommandController provides usecase.Command functionality to external drivers
type CommandController interface {
	Create(*usecase.CreateRequest) (usecase.EntityID, error)
}

// CommandControllerImpl implements CommandController
type CommandControllerImpl struct {
	service usecase.CommandService
}

var _ CommandController = &CommandControllerImpl{}

// NewCommandControllerImpl is a constructor
func NewCommandControllerImpl(service usecase.CommandService) *CommandControllerImpl {
	return &CommandControllerImpl{
		service: service,
	}
}

// Create creates new TODOs
func (cc *CommandControllerImpl) Create(request *usecase.CreateRequest) (usecase.EntityID, error) {
	id, err := cc.service.Create(request)
	if err != nil {
		return usecase.InvalidEntityID, convertError(err)
	}
	return id, nil
}

func convertError(in error) error {
	fieldError := &domain.FieldError{}
	if errors.As(in, &fieldError) {
		return fieldError
	}

	return in
}

// --------------- State ---------------

// DummyStateModifier implements usecase.StateModifier
type DummyStateModifier struct{}

var _ usecase.StateModifier = &DummyStateModifier{}

// Create is a dummy method
func (dsm *DummyStateModifier) Create(todo *domain.TODO) (usecase.EntityID, error) {
	return usecase.EntityID(101), nil
}
