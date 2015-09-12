package user

// Ensures persistance of user accounts in Postgresql database

import (
	"bleuvanille/config"
	"log"
)

// SaveUser inserts a user into the database
func Save(user User) (User, error) {
	_, err := config.Db().Query("INSERT INTO users (id, email, firstname, lastname, hash, isadmin, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7);", user.ID, user.Email, user.Firstname, user.Lastname, user.Hash, user.IsAdmin, user.CreatedAt)
	return user, err
}

// SavePassword updates the password of a given user into the database
func SavePassword(user *User, newPassword string) error {
	_, err := config.Db().Query("UPDATE users SET hash = $1 WHERE email = $2;", newPassword, user.Email)
	return err
}

// LoadByEmail returns the user object for a given email, if any
func LoadByEmail(email string) (*User, error) {
	var user User
	row := config.Db().QueryRow("SELECT * FROM users WHERE email = $1;", email)
	err := row.Scan(&user.ID, &user.Email, &user.Firstname, &user.Lastname, &user.Hash, &user.IsAdmin, &user.CreatedAt)

	return &user, err
}

// LoadByID returns the user object for a given user ID, if any
func LoadByID(ID string) (*User, error) {
	var user User
	row := config.Db().QueryRow("SELECT * FROM users WHERE id = $1;", ID)
	err := row.Scan(&user.ID, &user.Email, &user.Firstname, &user.Lastname, &user.Hash, &user.IsAdmin, &user.CreatedAt)

	return &user, err
}

// LoadAll returns the list of all users in the database
// for security concerns, the password hashes are not returned
// I don't think there is any case in which they are required
func LoadAll() (*Users, error) {
	var users Users

	rows, err := config.Db().Query("SELECT * FROM users")

	if err != nil {
		log.Printf("Cannot query user list: %s", err)
		return nil, err
	}

	// Close rows after all readed
	defer rows.Close()

	for rows.Next() {
		var user User

		err := rows.Scan(&user.Email, &user.Firstname, &user.Lastname, &user.CreatedAt)

		if err != nil {
			log.Println(err)
		}

		users = append(users, user)
	}

	return &users, err
}

// Delete delete the given user from the database
func Delete(user *User) error {
	_, err := config.Db().Query("DELETE FROM users WHERE email=$1;", user.Email)

	if err != nil {
		log.Printf("Cannot delete user: %s", err)
		return err
	}
	return nil
}
