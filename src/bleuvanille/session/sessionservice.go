package session

import (
	"bleuvanille/config"
	"errors"
	"fmt"

	ara "github.com/diegogub/aranGO"
)

// Save inserts a session into the database
func Save(sess *Session) error {
	errorMap := config.Context().Save(sess)
	if value, ok := errorMap["error"]; ok {
		errorString := fmt.Sprintf("Impossible to save the session with ID %v, userID %v, isAdmin %v, ExpiresAt %v : %s\n", sess.SessionID, sess.UserID, sess.IsAdmin, sess.ExpiresAt, value)
		return errors.New(errorString)
	}
	return nil
}

// GetByID returns the sess object for a given ID, if any,
// or nil if the session does not exist in the database
func GetByID(ID string) (*Session, error) {
	queryString := fmt.Sprintf("FOR s in sessions filter s.SessionID == %q RETURN s", ID)

	arangoQuery := ara.NewQuery(queryString)
	cursor, err := config.Db().Execute(arangoQuery)

	//return an error
	if err != nil {
		return nil, err
	}

	var result Session

	if cursor.Result != nil && len(cursor.Result) > 0 {
		cursor.FetchOne(&result)
		return &result, nil
	}
	return nil, nil
}

// Update updates the session in the database
// TODO we should renew the id is the admin privilege is given
// See https://www.owasp.org/index.php/Session_Management_Cheat_Sheet
func Update(sess *Session) error {
	errorMap := config.Context().Save(sess)

	if value, ok := errorMap["error"]; ok {
		errorString := fmt.Sprintf("Impossible to update the session with ID %v, userID %v, isAdmin %v, ExpiresAt %v : %s\n", sess.SessionID, sess.UserID, sess.IsAdmin, sess.ExpiresAt, value)
		return errors.New(errorString)
	}
	return nil
}

// Delete removes from the database the session for the given ID
func Delete(ID string) error {
	result, err := GetByID(ID)

	//return an error
	if err != nil {
		return err
	}

	if result == nil {
		errorString := fmt.Sprintf("Session with ID %s does not exists", ID)
		return errors.New(errorString)
	}

	errorMap := config.Context().Delete(result)

	if value, ok := errorMap["error"]; ok {
		errorString := fmt.Sprintf("Impossible to delete the session with ID %v : %s\n", result.SessionID, value)
		return errors.New(errorString)
	}

	return nil
}
