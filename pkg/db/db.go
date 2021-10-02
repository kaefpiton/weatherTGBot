package db

type Datastore interface {
	usersInterface
	stickersInterface
}

type usersInterface interface {
	InsertUser(usersFirstname, usersLastname string, chatid int64)error
}

type stickersInterface interface {
	GetStickersCodesByType(stickerTypeName string) ([]string, error)
}
