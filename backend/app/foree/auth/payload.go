package auth

type SignUpReq struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	ReferralCode string `json:"referralCode"`
}

type SessionReq struct {
	SessionId string `json:"sessionId"`
}

type VerifyEmailReq struct {
	SessionReq
	Code string `json:"code"`
}

type LoginReq struct {
	SessionReq
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ForgetPasswordUpdateReq struct {
	RetrieveCode string
	NewPassword  string
}
