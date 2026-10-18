package apis

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"myst/client/golang/core"
)

func (s *Server) LoginAccount(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if c.ShouldBindJSON(&req) != nil || req.Username == "" || req.Password == "" {
		Error(c, http.StatusBadRequest)
		return
	}
	token, err := s.app.LoginAccount(req.Username, req.Password)
	if errors.Is(err, core.ErrAuthenticationFailed) {
		Error(c, http.StatusUnauthorized)
		return
	}
	if err != nil {
		Error(c, http.StatusBadGateway, "account-server-unavailable")
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (s *Server) SignupAccount(c *gin.Context) {
	var req struct {
		Username       string `json:"username"`
		Password       string `json:"password"`
		MasterPassword string `json:"masterPassword"`
	}
	if c.ShouldBindJSON(&req) != nil {
		Error(c, http.StatusBadRequest)
		return
	}
	token, err := s.app.SignupAccount(req.Username, req.Password, req.MasterPassword)
	if errors.Is(err, core.ErrEnclaveExists) {
		Error(c, http.StatusConflict, "vault-already-registered")
		return
	}
	if err != nil {
		Error(c, http.StatusBadRequest, "registration-failed")
		return
	}
	c.JSON(http.StatusCreated, gin.H{"token": token})
}

func (s *Server) AccountStatus(c *gin.Context) {
	username, err := s.app.AccountStatus(c.GetHeader("X-Myst-Account"))
	if err != nil {
		Error(c, http.StatusUnauthorized)
		return
	}
	c.JSON(http.StatusOK, gin.H{"username": username})
}
