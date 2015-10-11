package session

import (
	"bleuvanille/config"
	"fmt"
	"log"
)

// Save inserts a session into the database
func Save(sess Session) error {

	fmt.Printf("Session  save: %v, %v, %v, %v\n", sess.SessionID, sess.UserID, sess.IsAdmin, sess.ExpiresAt)

	_, err := config.Db().Query("INSERT INTO sessions (id, user_id, is_admin, expires_at) VALUES ($1, $2, $3, $4);", sess.SessionID, sess.UserID, sess.IsAdmin, sess.ExpiresAt)
	if err != nil {
		log.Printf("Error: cannot save session with ID %v, userID %v, isAdmin %v, ExpiresAt %v\n", sess.SessionID, sess.UserID, sess.IsAdmin, sess.ExpiresAt)
	}
	return err
}

// GetByID returns the sess object for a given ID, if any,
// or nil if the session does not exist in the database
func GetByID(ID string) (*Session, error) {
	var sess Session
	row := config.Db().QueryRow("SELECT * FROM sessions WHERE id = $1;", ID)
	err := row.Scan(&sess.SessionID, &sess.UserID, &sess.IsAdmin, &sess.ExpiresAt)
	return &sess, err
}

// Update updates the session in the database
// TODO we should renew the id is the admin privilege is given
// See https://www.owasp.org/index.php/Session_Management_Cheat_Sheet
func Update(sess Session) error {
	_, err := config.Db().Query("UPDATE sessions SET is_admin = $2, expires_at = $3 WHERE id = $4;", sess.SessionID, sess.IsAdmin, sess.ExpiresAt)
	return err
}

// Delete removes from the database the session for the given ID
func Delete(ID string) error {
	_, err := config.Db().Query("DELETE FROM sessions WHERE id = $1;", ID)
	return err
}
