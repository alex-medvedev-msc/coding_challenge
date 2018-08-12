package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) GetAccounts(c *gin.Context) {
	accounts, err := s.accountRep.GetAccounts()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, accounts)
}

