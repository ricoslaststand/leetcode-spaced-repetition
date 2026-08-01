package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UserIDContextKey is the gin context key holding the authenticated user's ID. It is set by
// the auth middleware and must never be written by a handler.
const UserIDContextKey = "userID"

func FormatErrorBody(ctx *gin.Context, statusCode int, errorMessage string) {
	ctx.JSON(statusCode, gin.H{
		"error": errorMessage,
	})
}

// UserIDFromContext returns the authenticated user's ID, which the auth middleware placed in
// the context. The second return value is false when no middleware ran, which means the
// route was registered without auth — handlers must treat that as a failure rather than
// falling back to a default identity.
func UserIDFromContext(ctx *gin.Context) (uuid.UUID, bool) {
	value, exists := ctx.Get(UserIDContextKey)
	if !exists {
		return uuid.Nil, false
	}

	userID, ok := value.(uuid.UUID)
	if !ok || userID == uuid.Nil {
		return uuid.Nil, false
	}

	return userID, true
}
