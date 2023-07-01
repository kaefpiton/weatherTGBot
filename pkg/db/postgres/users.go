package postgres

import (
	"database/sql"
	"fmt"
	"time"
)

type User struct {
	ID        int64
	Firstname string
	Lastname  string
	ChatID    int64
}

func (db *DB) InsertUser(firstname, lastname string, chatID int64) error {
	userID, err := getUserIDByChatID(db, chatID)

	if err != nil {
		return err
	}

	if userID == -1 {
		//todo прикрутить лог
		fmt.Println("Создание нового пользователя")
		return createUser(db, firstname, lastname, chatID)
	} else {
		fmt.Println("Пользователь с chatID существует - меняем дату")
		return updateLastUsage(db, userID)
	}
}

func createUser(db *DB, firstname, lastname string, chatID int64) error {
	user := User{
		Firstname: firstname,
		Lastname:  lastname,
		ChatID:    chatID,
	}

	_, err := db.Exec("INSERT INTO users (firstname, lastname, chat_id) values ($1, $2, $3)",
		user.Firstname,
		user.Lastname,
		user.ChatID)

	return err
}

func updateLastUsage(db *DB, userID int64) error {
	_, err := db.Exec("UPDATE users SET last_usage = $1 WHERE id = $2",
		time.Now(),
		userID)

	return err
}

// todo возможно primary сделать chatID (посмотреть не повторяются ли)
func getUserIDByChatID(db *DB, chatID int64) (int64, error) {
	stmt, err := db.Prepare("SELECT ID FROM users WHERE chat_id = $1")
	if err != nil {
		return 0, err
	}

	var userID int64

	err = stmt.QueryRow(chatID).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, nil
		}

		return -1, err
	}
	return userID, nil
}
