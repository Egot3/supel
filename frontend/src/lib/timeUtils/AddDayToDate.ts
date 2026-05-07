export function AddDayToDate(date: Date, days: number): Date {
	date.setDate(date.getDate() + days);

	return date;
}
