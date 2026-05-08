export function GetWeekDateFromWeek(year: number, week: number): Date {
	const firstDayOfYear = new Date(Date.UTC(year, 0, 1));
	const daysOffset = (week - 1) * 7;

	const firstDayOfWeek = new Date(
		firstDayOfYear.setUTCDate(
			firstDayOfYear.getUTCDate() + daysOffset - firstDayOfYear.getUTCDay() + 1
		)
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
