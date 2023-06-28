package db

type Datastore interface {
	usersInterface
	stickersInterface
}

// todo вынести в usecase
type usersInterface interface {
	InsertUser(usersFirstname, usersLastname string, chatId int64) error
}

type stickersInterface interface {
	GetStickersCodesByType(stickerTypeName string) ([]string, error)
}
