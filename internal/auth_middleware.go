package internal

import (
	"net/http"

	"leetcode-spaced-repetition/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RemoteUserHeader is set by the authenticating reverse proxy (Authelia, via Traefik's
// ForwardAuth middleware) on every request it has already authenticated.
const RemoteUserHeader = "Remote-User"

// OwnerOnlyAuthMiddleware accepts requests the reverse proxy has authenticated as the single
// owner of this deployment and rejects everything else.
//
// Traefik's ForwardAuth middleware is the real gate — the API must not be reachable except
// through it. This check is defence in depth: if the middleware is ever detached, or the API
// port is published by mistake, requests arrive without a Remote-User header and are refused
// rather than silently served as the owner.
func OwnerOnlyAuthMiddleware(ownerUsername string, ownerUserID uuid.UUID) gin.HandlerFunc {
	return func(c *gin.Context) {
		remoteUser := c.GetHeader(RemoteUserHeader)

		// The empty check is deliberate and not redundant: if ownerUsername were ever
		// misconfigured to "", a request with no Remote-User header would compare equal
		// and this middleware would fail open.
		if remoteUser == "" || ownerUsername == "" || remoteUser != ownerUsername {
			utils.FormatErrorBody(c, http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		c.Set(utils.UserIDContextKey, ownerUserID)
		c.Next()
	}
}
