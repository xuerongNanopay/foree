type cause = "cause"

type SignInFormData = {
    email: string,
    password: string,
}

type SignInFormError = {
    [k in (keyof Partial<SignInFormData> | cause)]?: string;
}

type SignUpFormData = {
    email: string,
    password: string,
    rePassword: string,
}

type SignUpFormError = {
    [k in (keyof Partial<SignUpFormData> | cause)]?: string;
}

type VerifyEmailData = {
    code: string,
}

type VerifyEmailError = {
    [k in (keyof Partial<VerifyEmailData> | cause)]?: string;
}

type ForgetPasswordData = {
    email: string 
}

type ForgetPasswordError = {
    [k in (keyof Partial<SignUpFormData> | cause)]?: string;
}

type RenewPasswordData = {
    token: string,
    password: string,
    rePassword: string,
}

type RenewPasswordError = {
    [k in (keyof Partial<RenewPasswordData> | cause)]?: string;
}
