package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID int
	Name string
	Email string
	HashedPassword []byte
	Created time.Time
}

type UserModel struct {
	DB *sql.DB
}


func (m *User) Insert(name, email, password string) error {
	return nil
}

func (m *User) Authenticate(email, password string) (int, error) {
	return 0, nil
}

func (m *User) Exists(id int) (bool, error) {
	return false, nil
}
