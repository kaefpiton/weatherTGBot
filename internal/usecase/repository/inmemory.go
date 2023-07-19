package repository

// todo добавить дженерики
type Inmemory interface {
	Set(key, val any)
	Get(key any) any
}
