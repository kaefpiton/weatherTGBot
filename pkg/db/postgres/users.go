package postgres

import (
"database/sql"
"fmt"
"time"
)


type Users struct {
	Users_id 					int64
	Users_firstname				string
	Users_lastname				string
	Users_chatid				int64
	Users_date_of_last_usage	time.Time
}


func (db *DB)InsertUser(usersFirstname, usersLastname string, chatid int64)error{
	userID,err := getUserIDByChatID(db, chatid)

	if err!= nil{
		return err
	}

	if userID == 0{
		fmt.Println("Создание нового пользователя")
		return createUser(db, usersFirstname, usersLastname, chatid)
	} else{
		fmt.Println("Пользователь с chatid существует - меняем дату")
		return updateUserDateOfLastUsage(db, userID)
	}
}

func createUser(db*DB, usersFirstname, usersLastname string, chatid int64) error {
	user := Users{}

	user.Users_firstname = usersFirstname
	user.Users_lastname = usersLastname
	user.Users_chatid = chatid
	user.Users_date_of_last_usage = time.Now()

	_, err := db.Exec("INSERT INTO users (users_firstname, users_lastname,users_chatid, users_date_of_last_usage) values ($1, $2, $3, $4)",
		user.Users_firstname,
		user.Users_lastname,
		user.Users_chatid,
		user.Users_date_of_last_usage )
	return err
}
func updateUserDateOfLastUsage(db*DB, userID int64) error{

	_, err := db.Exec("UPDATE users SET users_date_of_last_usage = $1 WHERE users_id = $2",
		time.Now(),
		userID)

	return err
}
func getUserIDByChatID (db *DB, chatid int64) (int64, error)  {

	stmt, err := db.Prepare("SELECT users_id FROM users WHERE users_chatid = $1")
	if err != nil {
		return 0, err
	}

	var user Users

	err = stmt.QueryRow(chatid).Scan(&user.Users_id)
	if err != nil {
		if err == sql.ErrNoRows {
			//Обработка пустого результата
			return 0, nil
		}

		return 0, err
	}
	return chatid, nil
}
