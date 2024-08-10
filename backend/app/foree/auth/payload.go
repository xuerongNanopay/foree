package auth

type SignUpReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type VerifyEmailReq struct {
	Code string `json:"code"`
}
