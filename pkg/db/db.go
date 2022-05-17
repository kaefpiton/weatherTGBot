package db

type Datastore interface {
	usersInterface
	stickersInterface
}

type usersInterface interface {
	IsUserExist(usersFirstname string, chatid int64) (bool, error)
	CreateUser(usersFirstname, usersLastname string, chatid int64) error
	SetUserCity(city string, chatid int64) error
	UpdateUserDateOfLastUsage(chatid int64) error
	GetUserCity(chatid int64) (string, error)
}

type stickersInterface interface {
	GetStickersCodesByType(stickerTypeName string) ([]string, error)
}
