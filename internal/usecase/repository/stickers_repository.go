package repository

type UsersRepository interface {
	Create(firstname, lastname string, chatID int64) error
	UpdateLastUsage(chatID int64) error
	IsExist(chatID int64) bool
}
