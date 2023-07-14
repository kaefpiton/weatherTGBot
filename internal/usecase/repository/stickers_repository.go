package repository

type UsersRepository interface {
	Create(firstname, lastname, state string, chatID int64) error
	UpdateLastUsage(chatID int64) error
	GetUserStateByChatID(chatID int64) (string, error)
	SetUserState(chatID int64, state string) error
	IsExist(chatID int64) bool
}
