package amp

import (
	"net/http"
	"sync"
)

type Ctx struct {
	writer  http.ResponseWriter
	request *http.Request

	values   map[string]any
	valuesMu sync.Mutex
}
