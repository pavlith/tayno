package apis

import (
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"myst/client/golang/apis/generated"
	"myst/client/golang/core"
)

func sessionId(c *gin.Context) []byte {
	// GetCurrentUserID returns the username of the currently logged-in user
	sid, ok := c.Get("sessionId")
	if !ok {
		return nil
	}

	return sid.([]byte)
}

func (s *Server) Logout(c *gin.Context) {
	if err := s.app.LogoutAccount(c.GetHeader("X-Myst-Account")); err != nil {
		Error(c, http.StatusInternalServerError)
		return
	}
	if err := s.app.Logout(sessionId(c)); err != nil {
		Error(c, http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusNoContent)
}

func (s *Server) Authenticate(c *gin.Context) {
	var req generated.AuthenticateRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		Error(c, http.StatusBadRequest)
		return
	}

	sessionId, err := s.app.UnlockAccount(c.GetHeader("X-Myst-Account"), req.Password)
	if errors.Is(err, core.ErrInitializationRequired) {
		Error(c, http.StatusConflict)
		return
	} else if errors.Is(err, core.ErrAccountMismatch) {
		Error(c, http.StatusForbidden, "vault-account-mismatch")
		return
	} else if errors.Is(err, core.ErrAuthenticationRequired) {
		Error(c, http.StatusUnauthorized, "account-login-required")
		return
	} else if errors.Is(err, core.ErrAuthenticationFailed) {
		Error(c, http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	sid := base64.StdEncoding.EncodeToString(sessionId)

	c.JSON(http.StatusOK, sid)
}

func (s *Server) CurrentUser(c *gin.Context) {
	sid := sessionId(c)

	u, err := s.app.CurrentUser(sid)
	if errors.Is(err, core.ErrCredentialsNotFound) {
		Error(c, http.StatusNotFound)
		return
	} else if u == nil {
		Error(c, http.StatusUnauthorized)
		return
	}

	restUser, err := s.userToJSON(sid, *u)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, restUser)
}
