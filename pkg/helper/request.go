package helper

import "context"

var (
	UserNameContextKey = "user"
)

func GetUsernameFromContext(ctx context.Context) (string, bool) {
	username, ok := ctx.Value(UserNameContextKey).(string)
	return username, ok
}
