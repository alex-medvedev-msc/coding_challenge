package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetAccounts is an endpoint for GET /accounts which gets all accounts in system without any filtering
func (s *Server) GetAccounts(c *gin.Context) {
	accounts, err := s.accountRep.GetAccounts()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, accounts)
}

