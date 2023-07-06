package repository

type UsersRepository interface {
	InsertUser(userFirstname, usersLastname string, chatId int64) error
}
