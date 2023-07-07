package repository

type UsersRepository interface {
	InsertUser(userFirstname, usersLastname string, chatId int64) error
	CreateUser(firstname, lastname string, chatID int64) error
	UpdateLastUsage(chatID int64) error
	IsUserExist(chatID int64) bool
}
