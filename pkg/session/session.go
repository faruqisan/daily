package session

import (
	"context"
	"errors"
)

type typeKey string

var (
	k = typeKey("user_id")

	errUIDNotFound = errors.New("user id not found in context")

	// ErrNotAuthorized is represent error when user perform to do unauthorized action
	ErrNotAuthorized = errors.New("authorization failure, you are not allowed to perform this action")
)

// SetUIDToCtx function create new context based on given old context that has UID key and value
func SetUIDToCtx(old context.Context, userID int64) context.Context {
	return context.WithValue(old, k, userID)
}

// GetUIDFromCTX return user id from context
func GetUIDFromCTX(ctx context.Context) (int64, error) {

	if raw := ctx.Value(k); raw != nil {
		return raw.(int64), nil // strconv.ParseInt(raw.(string), 10, 64)
	}
	return 0, errUIDNotFound
}

// IsAuthorized function check if given user id same with user id exist on session context
func IsAuthorized(ctx context.Context, userID int64) error {
	uIDFromCTX, err := GetUIDFromCTX(ctx)
	if err != nil {
		return err
	}

	if uIDFromCTX != userID {
		return ErrNotAuthorized
	}

	return nil
}
