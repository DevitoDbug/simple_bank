package api

import (
	"github.com/gin-gonic/gin"
	db "simple_bank/db/sqlc"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	r := gin.Default()

	r.POST("/accounts", server.createAccount)

	return server
}


func errorResponse(err error) gin.H {
	//fix this thing to the things
	return gin.H("error": err.Error())
}
