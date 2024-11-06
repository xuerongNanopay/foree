import type { Actions } from './$types'
import { fail } from '@sveltejs/kit'

export const actions = {
	default: async ({request}) => {
		const data = await request.formData()
		const password = data.get('password') as string 
		const rePassword = data.get('rePassword') as string
		const token = data.get('token') as string

        const payload: RenewPasswordData = {
            password,
            rePassword,
            token,
        }

        console.log(password)

        if (payload.password !== payload.rePassword) {
            return fail<RenewPasswordError>(422, {
                rePassword: "password not match"
            })
        }
        //TODO:
	}
} satisfies Actions