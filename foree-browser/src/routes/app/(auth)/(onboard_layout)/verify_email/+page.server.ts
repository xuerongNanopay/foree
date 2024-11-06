import type { Actions } from './$types'

export const actions = {
	verify_code: async ({request}) => {
		const data = await request.formData()
		const email = data.get('emailVerifyCode') as string 
        //TODO:
        console.log("verify email", email)
	},
    resend_code: async({request}) => {
        console.log("resend code")
    }
} satisfies Actions