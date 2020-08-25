package session

import (
	"auth-rbac/internal/user"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type Session struct {
	ID         int64
	Token      string
	UserID     int64
	CreatedAt  time.Time
	ValidUntil time.Time
}

type Sessions interface {
	Upsert(session *Session) error
	FindByUserID(id int64) (*Session, error)
	FindByToken(token string) (*Session, error)
}

func CreateToken(id int64, login string) string {
	timeNow := time.Now()
	timeNowStr := timeNow.String()
	data := []byte(fmt.Sprintf("%v\n%s\n%s", id, login, timeNowStr))
	token := base64.StdEncoding.EncodeToString(data)

	return token
}

func DecodeToken(token string) (int64, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return 0, errors.Wrap(err, "can not decode token")
	}

	decodedSlice := strings.Split(string(decodedBytes), "\n")

	userID, err := strconv.ParseInt(decodedSlice[0], 10, 64)
	if err != nil {
		return 0, errors.Wrap(err, "failed to parse decode userID: %s")
	}

	return userID, nil
}

//CreateSes создают сессию
func CreateSes(user *user.User, duration string) (string, *Session, error) {
	var ses Session

	token := CreateToken(user.ID, user.Login)

	dur, err := time.ParseDuration(duration)
	if err != nil {
		return "", nil, errors.Wrap(err, "failed to parse decode userID: %s")
	}

	ses = Session{
		Token:      token,
		UserID:     user.ID,
		CreatedAt:  time.Now(),
		ValidUntil: time.Now().Add(dur),
	}

	return token, &ses, nil
}

func CheckValidSes(userToken string, s *Session) bool {
	const hour = 3 // fix time zone
	if s.Token == userToken && time.Now().Before(s.ValidUntil) {
		return true
	}

	return false
}
