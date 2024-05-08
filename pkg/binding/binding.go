package binding

import "net/http"

type Binder interface {
	Name() string
	Bind(*http.Request, any) error
}
