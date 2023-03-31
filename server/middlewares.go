package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Cors(ctx *gin.Context) {
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	ctx.Writer.Header().Set("Access-Control-Allow-Headers", "*")
	ctx.Writer.Header().Set("Access-Control-Allow-Methods", "*")

	if ctx.Request.Method == "OPTIONS" {
		ctx.Writer.WriteHeader(http.StatusOK)
		return
	}
}
