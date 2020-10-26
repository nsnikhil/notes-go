package handler

import (
	"net/http"
	"notes/pkg/http/internal/contract"
	"notes/pkg/http/internal/util"
	"notes/pkg/liberr"
	"notes/pkg/user"
)

type UserHandler struct {
	svc user.Service
}

func (uh *UserHandler) CreateUser(resp http.ResponseWriter, req *http.Request) error {
	var data contract.CreateUserRequest
	if err := util.ParseRequest(req, &data); err != nil {
		return liberr.WithArgs(liberr.Operation("UserHandler.CreateUser"), err)
	}

	//TODO: THINK IF THE VALIDATION SHOULD BE DELEGATED TO SVC LAYER ?

	_, err := uh.svc.CreateUser(data.Name, data.Email, data.Password)
	if err != nil {
		return liberr.WithArgs(liberr.Operation("UserHandler.CreateUser"), err)
	}

	//TODO: WRITE SUCCESS LOG
	util.WriteSuccessResponse(http.StatusCreated, contract.CreateUserResponse{Message: contract.UserCreationSuccess}, resp)
	return nil
}

func NewUserHandler(svc user.Service) *UserHandler {
	return &UserHandler{
		svc: svc,
	}
}
