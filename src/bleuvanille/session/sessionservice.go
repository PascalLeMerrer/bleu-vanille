package session

import (
	"bleuvanille/config"
)

// Save inserts a session into the database
func Save(session Session) (Session, error) {

	// Est-ce que je sauve l'utilsateur complet, ou seulement son Id, ou une partie de l'utilsateur (les champs fréquemment utilisés)
	_, err := config.Db().Query("INSERT INTO sessions (id, user_id, token, is_admin, expires_at) VALUES ($1, $2, $3, $4, $5);", session.ID, session.UserID, session.AuthToken, session.IsAdmin, session.ExpiresAt)
	return session, err
}

// GetByID returns the session object for a given ID, if any,
// or nil if the session does not exist in the database
func GetByID(ID string) (*Session, error) {
	var session Session
	row := config.Db().QueryRow("SELECT * FROM sessions WHERE id = $1;", ID)
	err := row.Scan(&session.ID, &session.UserID, session.AuthToken, &session.IsAdmin, &session.ExpiresAt)

	return &session, err
}

// Update updates the session in the database
func Update(session Session) error {
	_, err := config.Db().Query("UPDATE sessions SET token = $1, is_admin = $2, expires_at = $3 WHERE id = $4;", session.AuthToken, session.IsAdmin, session.ExpiresAt, session.ID)
	return err
}

// Delete removes from the database the session for the given ID
func Delete(ID string) error {
	_, err := config.Db().Query("DELETE FROM sessions WHERE id = $1;", ID)
	return err
}
