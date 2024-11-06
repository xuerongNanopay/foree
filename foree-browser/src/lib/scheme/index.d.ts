type SignInFormData = {
    email: string,
    password: string,
}

type SignInFormError = {
    [k in (keyof Partial<SignInFormData> | "cause")]?: string;
}

type SignUpFormData = {
    email: string,
    password: string,
    rePassword: string,
}

type SignUpFormError = {
    [k in (keyof Partial<SignUpFormData> | "cause")]?: string;
}