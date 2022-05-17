package postgres

import (
	"database/sql"
	"time"
)

type Users struct {
	Users_id                 int64
	Users_firstname          string
	Users_lastname           string
	Users_chatid             int64
	Users_date_of_last_usage time.Time
	Users_city               string
}

//PUBLIC
//проверка существования пользователя по имени и чат id
func (db *DB) IsUserExist(usersFirstname string, chatid int64) (bool, error) {
	var result = false

	stmt, err := db.Prepare("SELECT users_id FROM users WHERE users_firstname = $1 AND users_chatid = $2  ")
	if err != nil {
		return result, err
	}
	var user Users

	err = stmt.QueryRow(usersFirstname, chatid).Scan(&user.Users_id)
	if err != nil {
		if err == sql.ErrNoRows {
			//Обработка пустого результата
			return result, nil
		}
		return result, err
	}
	result = true
	return result, nil
}

//создает нового пользователя
func (db *DB) CreateUser(usersFirstname, usersLastname string, chatid int64) error {
	user := Users{}

	user.Users_firstname = usersFirstname
	user.Users_lastname = usersLastname
	user.Users_chatid = chatid
	user.Users_date_of_last_usage = time.Now()

	_, err := db.Exec("INSERT INTO users (users_firstname, users_lastname,users_chatid, users_date_of_last_usage) values ($1, $2, $3, $4)",
		user.Users_firstname,
		user.Users_lastname,
		user.Users_chatid,
		user.Users_date_of_last_usage)
	return err
}

//устанавливает город для полоьзователя
func (db *DB) SetUserCity(city string, chatid int64) error {
	userID, err := getUserIDByChatID(db, chatid)
	if err != nil {
		return err
	}

	_, err = db.Exec("UPDATE users SET users_city = $1 WHERE users_id = $2",
		city,
		userID)

	return err
}

//Обновляет дату входа пользователя
func (db *DB) UpdateUserDateOfLastUsage(chatid int64) error {
	userID, err := getUserIDByChatID(db, chatid)
	if err != nil {
		return err
	}

	_, err = db.Exec("UPDATE users SET users_date_of_last_usage = $1 WHERE users_id = $2",
		time.Now(),
		userID)

	return err
}

//возвращает город пользователя
func (db *DB) GetUserCity(chatid int64) (string, error) {
	row := db.QueryRow(`SELECT users_city FROM users where users_chatid = $1`, chatid)

	user := Users{}
	err := row.Scan(&user.Users_city)
	if err != nil {
		return "", err
	}

	return user.Users_city, nil
}

//PRIVATE
func getUserIDByChatID(db *DB, chatid int64) (int64, error) {
	var result int64

	stmt, err := db.Prepare("SELECT users_id FROM users WHERE users_chatid = $1")
	if err != nil {
		return 0, err
	}

	var user Users

	err = stmt.QueryRow(chatid).Scan(&user.Users_id)
	if err != nil {
		if err == sql.ErrNoRows {
			//Обработка пустого результата
			return result, nil
		}
		return 0, err
	}

	result = user.Users_id
	return result, nil
}
