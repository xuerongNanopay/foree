import type { RequestHandler } from './$types'
import { redirect } from '@sveltejs/kit';


export const POST: RequestHandler = ({ request, cookies }) => {
	throw redirect(302, "sign_in")
}