import type { PageServerLoad } from './$types'
import { redirect } from '@sveltejs/kit';

export const load: PageServerLoad = ({ }) => {
	redirect(302, "sign_in")
}