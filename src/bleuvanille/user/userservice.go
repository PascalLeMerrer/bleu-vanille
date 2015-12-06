package user

// Ensures persistance of user accounts in Postgresql database

import (
	"bleuvanille/config"
	"errors"
	"fmt"
	"math"
	"strconv"

	ara "github.com/diegogub/aranGO"
)

// Save inserts a user into the database
func Save(user *User) error {
	errorMap := config.Context().Save(user)
	if value, ok := errorMap["error"]; ok {
		return errors.New(value)
	}
	return nil
}

// SavePassword updates the password of a given user into the database
func SavePassword(user *User, newPassword string) error {
	user.Hash = newPassword
	return Save(user)
}

// LoadByEmail returns the user object for a given email, if any
func LoadByEmail(email string) (*User, error) {
	var result User

	col := config.GetCollection(&result)
	cursor, err := col.Example(map[string]interface{}{"email": email}, 0, 1)

	//return an error
	if err != nil {
		return nil, err
	}

	if cursor.Result != nil && len(cursor.Result) > 0 {
		cursor.FetchOne(&result)
		return &result, nil
	}

	return nil, nil
}

// LoadByID returns the user object for a given user ID, if any
func LoadByID(ID string) (*User, error) {
	querystr := fmt.Sprintf("FOR u in users filter u.id == %q RETURN u", ID)

	arangoquery := ara.NewQuery(querystr)
	cursor, err := config.Db().Execute(arangoquery)

	//return an error
	if err != nil {
		return nil, err
	}

	var result User

	if cursor.Result != nil && len(cursor.Result) > 0 {
		cursor.FetchOne(&result)
		return &result, nil
	}
	return nil, nil
}

// LoadAll returns the list of all users in the database
// sort defines the sorting property name
// order must be either ASC or DESC
// offset is the start index
// limit defines the max number of results to be returned
// returns an array of user, the total number of users, an error
func LoadAll(sort string, order string, offset int, limit int) ([]User, int, error) {
	limitString := ""
	if limit > 0 {
		limitString = " LIMIT " + strconv.Itoa(offset) + ", " + strconv.Itoa(limit)
	} else {
		limitString = " LIMIT " + strconv.Itoa(offset) + ", " + strconv.Itoa(math.MaxUint16)
	}
	queryString := "FOR u in users SORT u." + sort + " " + order + limitString + " RETURN u"
	arangoQuery := ara.NewQuery(queryString)
	cursor, err := config.Db().Execute(arangoQuery)
	fmt.Printf("\n*****************\ncursor.FullCount is %+v\n", cursor.FullCount())

	if err != nil {
		fmt.Println(err)
		return nil, 0, err
	}
	result := make([]User, len(cursor.Result))
	err = cursor.FetchBatch(&result)
	if err != nil {
		fmt.Println(err)
		return nil, 0, err
	}
	fmt.Printf("\ncursor.FullCount is %+v\n*****************\n", cursor.FullCount())
	return result, cursor.FullCount(), nil
}

// Delete delete the given user from the database
func Delete(user *User) error {

	if user == nil {
		return errors.New("Impossible to delete nil user")
	}

	err := config.Context().Delete(user)

	if err != nil && len(err) > 0 {
		errorstring := fmt.Sprintf("Impossible to delete User %s : %s", user.Key, err)
		return errors.New(errorstring)
	}

	return nil
}
