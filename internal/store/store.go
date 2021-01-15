package store

type Store interface {
	GetParam(param string, secret bool) (string, error)
	PutParam(param string, value string) error
}