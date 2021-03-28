package usecase

import (
	"fmt"
	"time"

	"github.com/liampulles/dump-and-dive/cmd/todo-service/internal/domain"
)

// EntityID identifies an entity within a collection.
type EntityID int32

// InvalidEntityID represents an entity that does not exist.
const InvalidEntityID = EntityID(-1)

// CreateRequest contains data necessary to create a TODO.
type CreateRequest struct {
	Name    string
	Details string
	Due     time.Time
}

// StateModifier modifies the state (persistence) of TODOs.
type StateModifier interface {
	Create(*domain.TODO) (EntityID, error)
}

// CommandService coordinates usecases which may change something in the system.
type CommandService interface {
	Create(*CreateRequest) (EntityID, error)
}

// CommandServiceImpl implements CommandService
type CommandServiceImpl struct {
	factory       domain.TODOFactory
	stateModifier StateModifier
}

var _ CommandService = &CommandServiceImpl{}

// NewCommandServiceImpl is a constructor
func NewCommandServiceImpl(
	factory domain.TODOFactory,
	stateModifier StateModifier,
) *CommandServiceImpl {
	return &CommandServiceImpl{
		factory:       factory,
		stateModifier: stateModifier,
	}
}

// Create creates and persists a TODO
func (cs *CommandServiceImpl) Create(request *CreateRequest) (EntityID, error) {
	todo, err := cs.factory.Create(request.Name, request.Details, request.Due)
	if err != nil {
		return InvalidEntityID, fmt.Errorf("todo factory error: %w", err)
	}

	id, err := cs.stateModifier.Create(todo)
	if err != nil {
		return InvalidEntityID, fmt.Errorf("state modifier error: %w", err)
	}

	return id, nil
}
