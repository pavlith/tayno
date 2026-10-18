package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"myst/server/apis/generated"
	"myst/server/core"
)

func (s *Server) CreateKeystore(c *gin.Context) {
	userId := CurrentUser(c)

	var req generated.CreateKeystoreRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		Error(c, http.StatusBadRequest)
		return
	}

	k, err := s.app.CreateKeystore(req.Name, userId, req.Payload)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, keystoreToJSON(k))
}

func (s *Server) Keystore(c *gin.Context) {
	userId := CurrentUser(c)
	keystoreId := c.Param("keystoreId")

	k, err := s.app.UserKeystore(userId, keystoreId)
	if errors.Is(err, core.ErrKeystoreNotFound) {
		Error(c, http.StatusNotFound)
		return
	} else if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, keystoreToJSON(k))
}

func (s *Server) UpdateKeystore(c *gin.Context) {
	userId := CurrentUser(c)
	keystoreId := c.Param("keystoreId")

	var req generated.UpdateKeystoreRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		Error(c, http.StatusBadRequest)
		return
	}

	k, err := s.app.UpdateKeystore(userId, keystoreId, core.UpdateKeystoreOptions{
		Name:    req.Name,
		Payload: req.Payload,
	})
	if errors.Is(err, core.ErrKeystoreNotFound) {
		Error(c, http.StatusNotFound)
		return
	} else if errors.Is(err, core.ErrForbidden) {
		Error(c, http.StatusForbidden)
		return
	} else if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, keystoreToJSON(k))
}

func (s *Server) Keystores(c *gin.Context) {
	userId := CurrentUser(c)

	ks, err := s.app.UserKeystores(userId)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	gen := []generated.Keystore{}
	for _, k := range ks {
		gen = append(gen, keystoreToJSON(k))
	}

	c.JSON(http.StatusOK, gen)
}

func (s *Server) DeleteKeystore(c *gin.Context) {
	userId := CurrentUser(c)
	keystoreId := c.Param("keystoreId")

	err := s.app.DeleteKeystore(userId, keystoreId)
	if errors.Is(err, core.ErrKeystoreNotFound) {
		Error(c, http.StatusNotFound)
		return
	} else if errors.Is(err, core.ErrForbidden) {
		Error(c, http.StatusForbidden)
		return
	} else if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
