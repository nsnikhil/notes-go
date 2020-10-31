package handler

import (
	"net/http"
	"notes/pkg/http/contract"
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
	_, err := uh.svc.Create(data.Name, data.Email, data.Password)
	if err != nil {
		return liberr.WithArgs(liberr.Operation("UserHandler.CreateUser"), err)
	}

	//TODO: WRITE SUCCESS LOG
	util.WriteSuccessResponse(http.StatusCreated, contract.CreateUserResponse{Message: contract.UserCreationSuccess}, resp)
	return nil
}

func (uh *UserHandler) LoginUser(resp http.ResponseWriter, req *http.Request) error {
	var data contract.LoginUserRequest
	if err := util.ParseRequest(req, &data); err != nil {
		return liberr.WithArgs(liberr.Operation("UserHandler.LoginUser"), err)
	}

	at, rf, err := uh.svc.Login(data.Email, data.Password)
	if err != nil {
		return liberr.WithArgs(liberr.Operation("UserHandler.LoginUser"), err)
	}

	util.WriteSuccessResponse(http.StatusOK, contract.LoginUserResponse{AccessToken: at, RefreshToken: rf}, resp)
	return nil
}

func NewUserHandler(svc user.Service) *UserHandler {
	return &UserHandler{
		svc: svc,
	}
}
