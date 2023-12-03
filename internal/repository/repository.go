package repository

import "errors"

var (
	ErrUserAlreadySubscribed = errors.New("user has already subscribed")
	ErrUserNotExists         = errors.New("user does not exists")
)
