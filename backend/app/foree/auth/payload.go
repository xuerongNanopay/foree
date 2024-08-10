package auth

type SignUpReq struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	ReferralCode string `json:"referralCode"`
}

type VerifyEmailReq struct {
	Code string `json:"code"`
}
