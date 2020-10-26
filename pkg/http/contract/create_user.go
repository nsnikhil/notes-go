package contract

const UserCreationSuccess = "user created successfully"

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	Message string `json:"message"`
}
