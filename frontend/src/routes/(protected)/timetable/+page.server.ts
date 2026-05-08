import type { Day, Timetable } from '$lib/types/timetable';
import type { PageServerLoad } from './$types';
import { fail, type Actions } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ /* fetch ,*/ cookies }) => {
	const token = cookies.get('auth_token');
	if (!token) {
		console.log('No token on load');
		return { currentWeek: NaN };
	}

	return { currentWeek: 1 };
};

export const actions: Actions = {
	getTimetable: async ({ /* fetch, */ cookies, request }) => {
		const token = cookies.get('auth_token');
		if (!token) {
			return fail(401, {
				err: 'NO_AUTH',
				errorMessage: 'You are not registered!'
			});
		}

		const data = await request.formData();
		const weekNumber = Number(data.get('weekNumber'));
		const dayNumber: Day = Number(data.get('dayNumber')!);
		const group = data.get('group');

		return { timetable: [{ day: dayNumber }] } as { timetable: Timetable };
	}
};
