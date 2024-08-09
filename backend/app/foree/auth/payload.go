package auth

type SignUpReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type VerifyEmailReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
