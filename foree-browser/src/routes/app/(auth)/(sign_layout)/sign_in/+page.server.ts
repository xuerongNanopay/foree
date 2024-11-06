import type { Actions } from './$types'
import { fail } from '@sveltejs/kit'

export const actions = {
	default: async ({request}) => {
		const data = await request.formData()
		const email = data.get('email') as string 
		const password = data.get('password') as string

        const errors = validateSignInForm({email, password})
        if (errors != null) {
            return fail<SignInFormError>(422, {
                ...errors,
                "cause": "TODO: error"
            })
        }
	}
} satisfies Actions

function validateSignInForm(data: SignInFormData): null|Partial<SignInFormData> {
    //TODO: real validation
    if (data.email === "aa@qq.com") {
        return {
            email: "error email"
        }
    }
    return null
}