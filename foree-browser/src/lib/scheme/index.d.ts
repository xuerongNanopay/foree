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
    referrerReference?: string,
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
    retrieveCode: string,
    password: string,
    rePassword: string,
}

type RenewPasswordError = {
    [k in (keyof Partial<RenewPasswordData> | cause)]?: string;
}

type CreateUserData = {
    firstName: string,
    middleName?: string,
    lastName: string,
    dob: string,
    pob: string,
    nationality: string,
    address1: string,
    address2?: string,
    city: string,
    province: string,
    country: string,
    postalCode: string,
    phoneNumber: string,
    identificationType: string,
    identificationValue: string,
    avatarUrl?: string
}

type CreateUserError = {
    [k in (keyof Partial<CreateUserData> | cause)]?: string;
}