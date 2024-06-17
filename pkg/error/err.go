// Github Repository: https://github.com/joseph-beck/amp
// GoDocs: https://pkg.go.dev/github.com/joseph-beck/amp

// Package Error is used for defining structured errors.
package error

type Error struct {
	Status  int
	Message string
}
