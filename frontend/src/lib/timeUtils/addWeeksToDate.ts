export function AddWeeksToDate(date: Date, weeks: number): Date {
	date.setDate(date.getDate() + weeks * 7);

	return date;
}
