package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type createAccountParams struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required",oneof=USD KSH`
}

func (s *Server) createAccount(ctx gin.Context) {
	var req createAccountParams
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

}
