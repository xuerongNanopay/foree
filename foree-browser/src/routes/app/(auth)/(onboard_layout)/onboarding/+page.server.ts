import type { Actions } from './$types'
import { fail } from '@sveltejs/kit'

export const actions = {
	default: async ({request}) => {
		const data = await request.formData()
		const firstName = data.get("firstName") as string 
		const rePassword = data.get('rePassword') as string
		const retrieveCode = data.get('retrieveCode') as string

        const payload = {
            firstName
        }

        console.log(payload)

	}
} satisfies Actions