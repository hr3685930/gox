package goo

// Error GooError
type Error interface {
	error
	GetStack() string
}