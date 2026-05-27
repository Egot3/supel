import { GetDayDateFromWeek, getWeekNumber } from '$lib/timeUtils/DateFromWeek';
import {
	type Lesson,
	type LessonCooked,
	type Period,
	type Subject,
	type TimetableEntry
} from '$lib/types/timetable';

import type { PageServerLoad } from './$types';
import { type Actions } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ fetch, cookies }) => {
	const token = cookies.get('auth_token');
	if (!token) {
		console.log('No token on load');
		return { currentWeek: 1, currentDay: 1 };
	}

	const date = new Date();
	const week = getWeekNumber(date);
	const day = date.getDay();

	const ownLearningGroupRequest = await fetch(`http://localhost/v1/group/self`, {
		headers: {
			Authorization: `Bearer ${token}`
		}
	});
	if (!ownLearningGroupRequest.ok) {
		return { currentWeek: week, currentDay: day };
	}
	const ownLearningGroups = await ownLearningGroupRequest.json();

	const Lessons: Lesson[] = [];
	ownLearningGroups.groups.array.forEach(async (group: { uuid: string }) => {
		console.log(group.uuid);
		const timetable = await fetch(`http://localhost/v1/timetable/current?${group.uuid}`, {
			headers: {
				Authorization: `Bearer ${token}`
			}
		});
		if (!timetable.ok) {
			if (timetable.status === 404) {
				return { currentWeek: week, currentDay: day, lessons: [] };
			}
			return { currentWeek: week, currentDay: day };
		}
		const timetableUUID: string = (await timetable.json()).uuid;

		const LessonsResponse = await fetch(
			`http://localhost/v1/timetable/lesson/timetable/${timetableUUID}`,
			{
				headers: {
					Authorization: `Bearer ${token}`
				}
			}
		);
		if (!LessonsResponse.ok) {
			return { currentWeek: week, currentDay: day };
		}

		Lessons.push((await LessonsResponse.json()) as Lesson);
	});

	const lessonPromises = Lessons.map(async (lesson: Lesson) => {
		const TimetableEntryRequest = await fetch(
			`http://localhost/v1/timetable/entry/${lesson.TimetableEntryUUID}`,
			{
				headers: {
					Authorization: `Bearer ${token}`
				}
			}
		);
		if (!TimetableEntryRequest.ok) {
			console.log(TimetableEntryRequest);
			throw { currentWeek: week, currentDay: day, lessons: [] };
		}

		const TimetableEntry: TimetableEntry = await TimetableEntryRequest.json();

		const PeriodResponse = await fetch(
			`http://localhost/v1/timetable/period/${TimetableEntry.periodUUID}`,
			{
				headers: {
					Authorization: `Bearer ${token}`
				}
			}
		);
		if (!PeriodResponse.ok) {
			console.log(PeriodResponse);
			throw { currentWeek: week, currentDay: day, lessons: [] };
		}

		const Period: Period = await PeriodResponse.json();

		const SubjectResponse = await fetch(`/v1/timetable/subject/${TimetableEntry.subjectUUID}`, {
			headers: {
				Authorization: `Bearer ${token}`
			}
		});
		if (!SubjectResponse.ok) {
			console.log(SubjectResponse);
			throw { currentWeek: week, currentDay: day, lessons: [] };
		}

		const Subject: Subject = await SubjectResponse.json();

		return {
			UUID: lesson.UUID,
			TimetableEntryUUID: TimetableEntry.UUID,
			Date: lesson.Date,
			Cancelled: lesson.Cancelled,
			Name: Subject.Name,
			Position: Period.Position,
			StartTime: Period.StartTime,
			EndTime: Period.EndTime
		} as LessonCooked;
	});

	const lessonsCooked = await Promise.all(lessonPromises);
	return { currentWeek: week, currentDay: day, lessons: lessonsCooked };
};

export const actions: Actions = {
	getTimetable: async ({ fetch, cookies }) => {
		const token = cookies.get('auth_token');
		if (!token) {
			console.log('No token on load');
			return { currentWeek: 1, currentDay: 1 };
		}

		const date = new Date();
		const week = getWeekNumber(date);
		const day = date.getDay();
		const ownLearningGroupRequest = await fetch(`http://localhost/v1/group/self`, {
			headers: {
				Authorization: `Bearer ${token}`
			}
		});
		if (!ownLearningGroupRequest.ok) {
			return { currentWeek: week, currentDay: day };
		}
		const ownLearningGroups = await ownLearningGroupRequest.json();

		const Lessons: Lesson[] = [];
		ownLearningGroups.groups.array.forEach(async (group: { uuid: string }) => {
			console.log(group.uuid);
			const timetable = await fetch(
				`http://localhost/v1/timetable/at?${GetDayDateFromWeek(2026, week, day).toString()}`,
				{
					headers: {
						Authorization: `Bearer ${token}`
					}
				}
			);
			if (!timetable.ok) {
				if (timetable.status === 404) {
					return { currentWeek: week, currentDay: day, lessons: [] };
				}
				return { currentWeek: week, currentDay: day };
			}
			const timetableUUID: string = (await timetable.json()).uuid;

			const LessonsResponse = await fetch(
				`http://localhost/v1/timetable/lesson/timetable/${timetableUUID}`,
				{
					headers: {
						Authorization: `Bearer ${token}`
					}
				}
			);
			if (!LessonsResponse.ok) {
				return { currentWeek: week, currentDay: day };
			}

			Lessons.push((await LessonsResponse.json()) as Lesson);
		});

		const lessonPromises = Lessons.map(async (lesson: Lesson) => {
			const TimetableEntryRequest = await fetch(
				`http://localhost/v1/timetable/entry/${lesson.TimetableEntryUUID}`,
				{
					headers: {
						Authorization: `Bearer ${token}`
					}
				}
			);
			if (!TimetableEntryRequest.ok) {
				throw { currentWeek: week, currentDay: day, lessons: [] };
			}

			const TimetableEntry: TimetableEntry = await TimetableEntryRequest.json();

			const PeriodResponse = await fetch(
				`http://localhost/v1/timetable/period/${TimetableEntry.periodUUID}`,
				{
					headers: {
						Authorization: `Bearer ${token}`
					}
				}
			);
			if (!PeriodResponse.ok) {
				console.log(PeriodResponse);
				throw { currentWeek: week, currentDay: day, lessons: [] };
			}

			const Period: Period = await PeriodResponse.json();

			const SubjectResponse = await fetch(`/v1/timetable/subject/${TimetableEntry.subjectUUID}`, {
				headers: {
					Authorization: `Bearer ${token}`
				}
			});
			if (!SubjectResponse.ok) {
				console.log(SubjectResponse);
				throw { currentWeek: week, currentDay: day, lessons: [] };
			}

			const Subject: Subject = await SubjectResponse.json();

			return {
				UUID: lesson.UUID,
				TimetableEntryUUID: TimetableEntry.UUID,
				Date: lesson.Date,
				Cancelled: lesson.Cancelled,
				Name: Subject.Name,
				Position: Period.Position,
				StartTime: Period.StartTime,
				EndTime: Period.EndTime
			} as LessonCooked;
		});

		const lessonsCooked = await Promise.all(lessonPromises);
		return { currentWeek: week, currentDay: day, lessons: lessonsCooked };
	}
};
