package handler

import (
	"net/http"
	"notes/pkg/http/internal/util"
)

func PingHandler() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		util.WriteSuccessResponse(http.StatusOK, "pong", resp)
	}
}
