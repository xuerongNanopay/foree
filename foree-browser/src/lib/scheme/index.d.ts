type SignUpFormData = {
    email: string,
    password: string
}

type SignUpFormError = {
    [k in (keyof Partial<SignUpFormData> | "cause")]?: string;
}