package repository

type Inmemory interface {
	Set(key, val any)
	Get(key any) any
}
