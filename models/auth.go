package models

var (
	// TokenExpirationInMinutes : 30min
	TokenExpirationInMinutes = 30
)

// UserCredentials : Models the structure of user credentials in login request body
type UserCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserInfos : User Infos from redis
type UserInfos struct {
	UserID string `json:"userID"`
}

// ChangePasswordRequest : Change password request body
type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}
