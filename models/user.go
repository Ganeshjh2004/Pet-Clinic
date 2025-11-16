package models

import (
	"database/sql"
	"petclinic/db"
	"petclinic/utils"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"` // store hashed passwords in production!
	Role     string `json:"role"`     // "owner" or "staff"
}

// GetUserByEmail fetches user data by email, used for login
func GetUserByEmail(email string) (*User, error) {
	var u User
	err := db.DB.QueryRow("SELECT id, email, password, role FROM users WHERE email=$1", email).
		Scan(&u.ID, &u.Email, &u.Password, &u.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.Warn("No user found with email: %s", email)
			return nil, nil
		}
		utils.Error("GetUserByEmail error: %v", err)
		return nil, err
	}
	return &u, nil
}
