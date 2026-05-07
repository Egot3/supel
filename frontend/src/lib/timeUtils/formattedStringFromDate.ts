export function FormattedStringFromDate(date: Date): string {
	const [yyyy, mm, dd] = date.toISOString().split('T')[0].split('-');
	return `${dd}/${mm}/${yyyy}`;
}
