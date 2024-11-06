import type { Actions } from './$types'

export const actions = {
	default: async ({request}) => {
		const data = await request.formData()
		const email = data.get('email') as string 
        //TODO:
        console.log("forget_password", email)
	}
} satisfies Actions