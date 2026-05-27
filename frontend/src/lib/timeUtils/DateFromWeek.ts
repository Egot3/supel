export function GetWeekDateFromWeek(year: number, week: number): Date {
	const firstDayOfYear = new Date(Date.UTC(year, 0, 1));
	const daysOffset = (week - 1) * 7;

	const firstDayOfWeek = new Date(
		firstDayOfYear.setUTCDate(firstDayOfYear.getUTCDate() + daysOffset - firstDayOfYear.getUTCDay())
	);

	return firstDayOfWeek;
}

export function GetDayDateFromWeek(year: number, week: number, day: number): Date {
	const firstDayOfYear = new Date(Date.UTC(year, 0, 1));
	const daysOffset = (week - 1) * 7 + day;

	const firstDayOfWeek = new Date(
		firstDayOfYear.setUTCDate(
			firstDayOfYear.getUTCDate() + daysOffset - firstDayOfYear.getUTCDay() + 1
		)
	);

	return firstDayOfWeek;
}

export const getWeekNumber = (date: Date): number => {
	const firstDayOfYear = new Date(date.getFullYear(), 0, 1);
	const pastDaysOfYear = (date.getTime() - firstDayOfYear.getTime()) / 86400000;
	return Math.ceil((pastDaysOfYear + firstDayOfYear.getDay() + 1) / 7);
};
