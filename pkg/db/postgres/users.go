package postgres

import (
	"database/sql"
	"fmt"
	"time"
)

type Users struct {
	ID        int64
	Firstname string
	Lastname  string
	ChatID    int64
}

func (db *DB) InsertUser(usersFirstname, usersLastname string, chatid int64) error {
	userID, err := getUserIDByChatID(db, chatid)

	if err != nil {
		return err
	}

	if userID == -1 {
		//todo прикрутить лог
		fmt.Println("Создание нового пользователя")
		return createUser(db, usersFirstname, usersLastname, chatid)
	} else {
		fmt.Println("Пользователь с chatID существует - меняем дату")
		return updateLastUsage(db, userID)
	}
}

func createUser(db *DB, usersFirstname, usersLastname string, chatID int64) error {
	var user Users

	user.Firstname = usersFirstname
	user.Lastname = usersLastname
	user.ChatID = chatID

	_, err := db.Exec("INSERT INTO users (firstname, lastname, chat_id) values ($1, $2, $3)",
		user.Firstname,
		user.Lastname,
		user.ChatID)
	return err
}

func updateLastUsage(db *DB, userID int64) error {
	_, err := db.Exec("UPDATE users SET last_usage = $1 WHERE chat_id = $2",
		time.Now(),
		userID)

	return err
}
func getUserIDByChatID(db *DB, chatID int64) (int64, error) {
	stmt, err := db.Prepare("SELECT ID FROM users WHERE chat_id = $1")
	if err != nil {
		return 0, err
	}

	var user Users

	err = stmt.QueryRow(chatID).Scan(&user.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, nil
		}

		return -1, err
	}
	return chatID, nil
}
