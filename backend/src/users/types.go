package users

type UserDataPayload struct {
	Email string `json:"email"`
	Hashed_password string `json:"hashedPassword"`
}