package armadan

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

const DEFAULT_TIME_FORMAT = "2006-01-02 15:04:05"

type Validator interface {
	ValidateIdParam(*http.Request, string) (string, error)
	Validate(*http.Request, any) error
}

func MapEntities[E, M any](entities []E, mapFn func(any) *M) []M {
	models := make([]M, len(entities))

	for i, entity := range entities {
		models[i] = *mapFn(entity)
	}

	return models
}

func GetId() string {
	return uuid.NewString()
}

func ParseTime(val string) time.Time {
	parsed, err := time.Parse(DEFAULT_TIME_FORMAT, val)

	if err != nil {
		return time.Now()
	}

	return parsed
}

func ToSqlBool(val bool) int64 {
	if val {
		return 1
	}

	return 0
}
