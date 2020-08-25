package user

import (
	"crypto/md5"
	"encoding/hex"
	"time"

	"github.com/pkg/errors"
)

// User сущность юзера
type User struct {
	ID        int64     `json:"-"`
	Login     string    `json:"login"`
	Password  string    `json:"pass"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"-"`
}

type Users interface {
	Create(u *User) error
	FindUserByID(id int64) (*User, error)
	FindUserByLogin(login string) (*User, error)
}

func HashPass(pass string) (string, error) {
	hasher := md5.New() // nolint
	_, err := hasher.Write([]byte(pass))

	if err != nil {
		return "", errors.Wrap(err, "hash failed %s")
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func CheckValidUser(user *User) error {
	if user.Password == "" {
		err := errors.New("enter password")

		return err
	}

	if user.Login == "" {
		err := errors.New("enter login")

		return err
	}

	return nil
}
