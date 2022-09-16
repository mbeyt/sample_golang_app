package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const (
	appSalt     = ""
	appHomepage = ""
)

type SsoRequest struct {
	ResourceUUID string `param:"resource_uuid" form:"resource_uuid"`
	Token        string `param:"token" form:"token"`
	Timestamp    string `param:"timestamp" form:"timestamp"`
	Email        string `param:"user_email" form:"user_email"`
	Id           string `param:"user_id" form:"user_id"`
}

type UnauthorizedError struct{}

func (e *UnauthorizedError) Error() string {
	return "Unable to validate token"
}

func (s *server) sso(req *SsoRequest) (*http.Cookie, error) {
	authorized, err := validToken(req.Token, req.Timestamp, req.ResourceUUID)
	if err != nil {
		return nil, err
	} else if !authorized {
		return nil, &UnauthorizedError{}
	}

	expiration := time.Now().Add(365 * 24 * time.Hour)
	//TODO: this probably also needs some sort of auth token cookie
	cookie := http.Cookie{Name: "uuid", Value: req.ResourceUUID, Expires: expiration}
	return &cookie, nil
}

func validToken(token string, timestamp string, uuid string) (bool, error) {
	// has this timestamp expired?
	i, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return false, err
	}
	tm := time.Unix(i, 0)
	if time.Since(tm).Minutes() > 2 {
		return false, nil
	}

	// is this token valid?
	decodedToken, err := hex.DecodeString(token)
	if err != nil {
		return false, err
	}
	message := []byte(fmt.Sprintf("%s:%s", timestamp, uuid))

	hash := hmac.New(sha256.New, []byte(appSalt))
	hash.Write(message)

	return hmac.Equal(hash.Sum(nil), []byte(decodedToken)), nil
}
