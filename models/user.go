package models

import (
	"ems_backend_go/db"
	"ems_backend_go/utils"
	"errors"
)

type User struct {
	ID       int64
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

func (u User) Save() error {
	query := `INSERT INTO users(email,password) VALUES(?,?)`
	stmnt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmnt.Close()
	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}
	res, err := stmnt.Exec(u.Email, hashedPassword)
	if err != nil {
		return err
	}
	userID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	u.ID = userID
	return nil
}

func GetUsers() ([]User, error) {
	query := `SELECT id, email, password FROM users`
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Email, &user.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil

}

func (u User) ValidateCredentials() error {
	query := `SELECT password FROM users WHERE email=?`
	row := db.DB.QueryRow(query, u.Email)
	var retrievedPassword string
	err := row.Scan(&retrievedPassword)
	if err != nil {
		return err
	}
	passwordIsValid := utils.CheckHashedPassword(u.Password, retrievedPassword)
	if passwordIsValid == true {
		return nil
	}
	return errors.New("Invalid Credentials")
}

func FetchIDByEmail(email string) (int64, error) {
	var id int64
	query := `SELECT id FROM users WHERE email=?`
	row := db.DB.QueryRow(query, email)
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
