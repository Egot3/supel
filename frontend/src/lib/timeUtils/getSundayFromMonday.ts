export function GetSundayFromMonday(date: Date): Date {
	date.setDate(date.getDate() + 6);

	return date;
}
