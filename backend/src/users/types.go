package users

type UserAuthPayload struct {
	Name string `json:"name"`
	Email string `json:"email"`
	Hashed_password string `json:"hashedPassword"`
}

type UserDataPayload struct {
	UID int32 `json:"uid"`
	Name string `json:"name"`
	Email string `json:"email"`
	Instagram_id string `json:"instagramId"`
	Verified bool `json:"verified"`
}

