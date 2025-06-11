package models

import "ems_backend_go/db"

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
	res, err := stmnt.Exec(u.Email, u.Password)
	if err != nil {
		return err
	}
	userID, err := res.LastInsertId()
	if err != nil {
		return errs
	}
	u.ID = userID
	return nil
}
