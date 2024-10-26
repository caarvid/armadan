package armadan

import (
	"net/http"

	"github.com/google/uuid"
)

type Validator interface {
	ValidateIdParam(*http.Request) (*uuid.UUID, error)
	Validate(*http.Request, interface{}) error
}

func MapEntities[E, M interface{}](entities []E, mapFn func(any) *M) []M {
	models := make([]M, len(entities))

	for i, entity := range entities {
		models[i] = *mapFn(entity)
	}

	return models
}
