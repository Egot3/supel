package types

type Verb string

const (
	UNDEFINED = iota
	GET       = "GET"
	POST      = "POST"
	PUT       = "PUT"
	DELETE    = "DELETE"
)
