package user

import (
	"errors"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/speps/go-hashids"
	"golang.org/x/crypto/bcrypt"
)

// User is any type of registered user
type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Hash      string
	AuthToken string `json:"token"`
	IsAdmin   bool
	CreatedAt time.Time `json:"createdAt"`
}

// Users is a list of User
type Users []User

// NewUser creates a User instance
func NewUser(email string, firstname string, lastname string, password string) (User, error) {
	var user User
	if email == "" {
		return user, errors.New("Cannot create user, email is missing")
	}
	if password == "" {
		return user, errors.New("Cannot create user, password is missing")
	}
	user.ID = generateID()
	user.Email = email
	user.Firstname = firstname
	user.Lastname = lastname
	user.IsAdmin = false
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return user, err
	}
	user.Hash = string(hash)
	user.CreatedAt = time.Now()

	return user, nil
}

// returns an 8 characters random string safe for use in URL
// the characters may be any of [A-Z][a-z][0-9]
func generateID() string {
	hashIDConfig := hashids.NewData()
	hashIDConfig.Salt = "zs4e6f80KDla1-2xcCD!34%<?23POsd"
	hashIDConfig.MinLength = 8
	hashIDConfig.Alphabet = hashids.DefaultAlphabet
	hash := hashids.NewWithData(hashIDConfig)

	randomInt := rand.New(rand.NewSource(time.Now().UnixNano())).Int63()
	intArray := intToIntArray(randomInt, 8)
	result, _ := hash.Encode(intArray)

	return result
}

// converts an int64 number to a fixed length array of int
func intToIntArray(value int64, length int) []int {
	result := make([]int, length)
	valueAsString := strconv.FormatInt(value, 10)

	fragmentLength := len(valueAsString) / length

	var startIndex, endIndex int
	var intValue int64
	var err error

	for index := 0; index < length; index++ {

		startIndex = index * fragmentLength
		endIndex = ((index + 1) * fragmentLength)

		if endIndex <= len(valueAsString) {
			intValue, err = strconv.ParseInt(valueAsString[startIndex:endIndex], 10, 0)
		} else {
			intValue, err = strconv.ParseInt(valueAsString[startIndex:], 10, 0)
		}

		if err != nil {
			log.Panicf("Error while converting string to int array %s", err)
		}
		result[index] = int(intValue)
	}
	return result
}
