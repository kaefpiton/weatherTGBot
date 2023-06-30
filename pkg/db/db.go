package db

type TgBotRepo interface {
	usersStore
	stickersStore
}

// todo вынести в usecase
type usersStore interface {
	InsertUser(userFirstname, usersLastname string, chatId int64) error
}

type stickersStore interface {
	GetStickersCodesByType(stickerTypeName string) ([]string, error)
}
