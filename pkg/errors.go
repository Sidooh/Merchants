package pkg

import "errors"

var (
	ErrInvalidMerchant = errors.New("merchant details are invalid")

	ErrInvalidUser = errors.New("user details are invalid")

	ErrInvalidAccount = errors.New("account details are invalid")

	ErrInvalidChannel = errors.New("channel is not supported")

	ErrUnauthorized = errors.New("unauthorized")

	ErrUnauthorizedMfa = errors.New("missing 2FA")

	ErrServerError = errors.New("something went wrong")
)
