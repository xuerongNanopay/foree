import type { Actions } from './$types'
import { fail } from '@sveltejs/kit'

export const actions = {
	sign_in: async ({request}) => {
		const data = await request.formData()
		const email = data.get('email') as string 
		const password = data.get('password') as string

        const errors = validateSignUpForm({email, password})
        if (errors != null) {
            return fail(422, {
                ...errors,
                "cause": "TODO: error"
            })
        }
	}
} satisfies Actions

function validateSignUpForm(data: SignUpFormData): null|Partial<SignUpFormData> {
    //TODO: real validation
    if (data.email === "aa@qq.com") {
        return {
            email: "error email"
        }
    }
    return null
}