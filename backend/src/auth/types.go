package auth

type VerificationPayload struct {
	Uid int32 `json:"uid"`
	Code string `json:"code"`
}
