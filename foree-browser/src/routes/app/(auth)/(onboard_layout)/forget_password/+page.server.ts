import type { Actions } from './$types'

export const actions = {
	forget_password: async ({request}) => {
		const data = await request.formData()
		const email = data.get('email') as string 
        //TODO:
        console.log("forget_password", email)
	}
} satisfies Actions