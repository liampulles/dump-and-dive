package domain

import (
	"fmt"
	"strings"
	"time"
)

// --------------- Errors ---------------

// FieldError is returned when field(s) fail validation checks.
type FieldError map[string][]string

var _ error = FieldError(nil)

// Error implements error
func (fe FieldError) Error() string {
	var items []string
	for field, issues := range fe {
		fieldIssue := fmt.Sprintf("%s=%s", field, strings.Join(issues, "|"))
		items = append(items, fieldIssue)
	}
	return fmt.Sprintf("field error(s): [%s]", strings.Join(items, ","))
}

// --------------- Basic services ---------------

// TimeProvider provides non-deterministic time data
type TimeProvider interface {
	Now() time.Time
}

// GoTimeProvider implements TimeProvider via the Go standard library
type GoTimeProvider struct{}

var _ TimeProvider = &GoTimeProvider{}

// NewGoTimeProvider is a constructor
func NewGoTimeProvider() *GoTimeProvider {
	return &GoTimeProvider{}
}

// Now returns the current local time
func (rp *GoTimeProvider) Now() time.Time {
	return time.Now()
}

// --------------- TODO ---------------

// TODO details a task or reminder that has an associated due date.
type TODO struct {
	name    string
	details string
	due     time.Time
}

// Name returns the name of the TODO
func (t *TODO) Name() string {
	return t.name
}

// Details returns additional context of the TODO.
// May be empty.
func (t *TODO) Details() string {
	return t.details
}

// Due indicates when the TODO needs to be completed by.
func (t *TODO) Due() time.Time {
	return t.due
}

// TODOFactory creates TODOs
type TODOFactory interface {
	Create(name string, details string, due time.Time) (*TODO, error)
}

// TODOFactoryImpl implements TODOFactory
type TODOFactoryImpl struct {
	timeProvider TimeProvider
}

var _ TODOFactory = &TODOFactoryImpl{}

// NewTODOFactoryImpl is a constructor
func NewTODOFactoryImpl(timeProvider TimeProvider) *TODOFactoryImpl {
	return &TODOFactoryImpl{
		timeProvider: timeProvider,
	}
}

// Create creates TODOs
func (tf *TODOFactoryImpl) Create(
	name string,
	details string,
	due time.Time,
) (*TODO, error) {

	validErr := validate().
		notBlank("name", name).
		isInFuture("due", due, tf.timeProvider).
		build()
	if validErr != nil {
		return nil, fmt.Errorf("validation error: %w", validErr)
	}

	return &TODO{
		name:    name,
		details: details,
		due:     due,
	}, nil
}

// --------------- Validation ---------------

type fieldIssue struct {
	field string
	issue string
}

type validationBuilder struct {
	fieldIssues []*fieldIssue
}

func validate() *validationBuilder {
	return &validationBuilder{}
}

func (vb *validationBuilder) notBlank(field string, value string) *validationBuilder {
	if value != "" {
		return vb
	}
	return vb.logIssue(field, "is blank")
}

func (vb *validationBuilder) isInFuture(
	field string,
	value time.Time,
	timeProvider TimeProvider,
) *validationBuilder {
	if value.Unix() <= timeProvider.Now().Unix() {
		return vb
	}
	return vb.logIssue(field, "is not in the future")
}

func (vb *validationBuilder) logIssue(field string, issue string) *validationBuilder {
	vb.fieldIssues = append(vb.fieldIssues, &fieldIssue{
		field: field,
		issue: issue,
	})
	return vb
}

func (vb *validationBuilder) build() FieldError {
	if len(vb.fieldIssues) == 0 {
		return nil
	}

	fieldIssueMap := make(map[string][]string)
	for _, fieldIssue := range vb.fieldIssues {
		fieldIssueMap[fieldIssue.field] =
			append(fieldIssueMap[fieldIssue.field], fieldIssue.issue)
	}
	return FieldError(fieldIssueMap)
}
