<script lang="ts">
	import {
		GetWeekDateFromWeek,
		GetSundayFromMonday,
		FormattedStringFromDate,
		GetDayDateFromWeek
	} from '$lib/timeUtils';

	import { Button } from 'flowbite-svelte';

	let weekNumber = $state(1);
	let currentYear = $state(1970);
	let chosen = $state('');
	let weekDate = $derived.by((): string => {
		const dateMonday = GetWeekDateFromWeek(currentYear, weekNumber);
		const formattedDateMonday = FormattedStringFromDate(dateMonday);
		console.log(dateMonday);

		const dateSunday = GetSundayFromMonday(dateMonday);
		const formattedDateSunday = FormattedStringFromDate(dateSunday);

		return `${formattedDateMonday}->${formattedDateSunday}`;
	});
	let weekDaysDates: Record<string, string> = $derived.by(() => {
		const dayNames = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday'];

		const result: Record<string, string> = {};

		dayNames.forEach((day, idx) => {
			result[day] = FormattedStringFromDate(GetDayDateFromWeek(currentYear, weekNumber, idx));
		});
		console.log(result);

		return result;
	});
</script>

<div
	class="col-start-2 row-start-1 col-end-12 row-end-3 bg-dark text-[min(clamp(1rem,10vh,7rem),clamp(1rem,10vw,7rem))] inline-grid items-center justify-center italic font-hollow tracking-[0.5rem] select-none pointer-events-none"
>
	TIMETABLE
</div>

<div class="col-start-5 col-end-9 row-start-3 row-end-3 bg-dark -mb-8 rounded-t-4xl">
	{weekDate}
</div>

<div
	class="
    col-start-1 row-start-4 row-end-12 col-end-13 bg-dark
    grid grid-cols-7 grid-rows-7 h-full w-full p-gutter
    rounded-4xl
    "
>
	<Button
		outline={chosen != 'Monday'}
		onclick={() => {
			chosen = 'Monday';
		}}
		class={`row-start-1 row-end-2 col-start-1 rounded-b-none rounded-r-none rounded-l-2xl ${chosen == 'Monday' ? 'bg-accent' : ''}`}
		>Monday</Button
	>
	<Button
		outline={chosen != 'Tuesday'}
		onclick={() => {
			chosen = 'Tuesday';
		}}
		class={`row-start-1 row-end-2 col-start-2 rounded-b-none rounded-r-none rounded-l-none ${chosen == 'Tuesday' ? 'bg-accent' : ''}`}
		>Tuesday</Button
	>
	<Button
		outline={chosen != 'Wednesday'}
		onclick={() => {
			chosen = 'Wednesday';
		}}
		class={`row-start-1 row-end-2 col-start-3 rounded-b-none rounded-r-none rounded-l-none ${chosen == 'Wednesday' ? 'bg-accent' : ''}`}
		>Wednesday</Button
	>
	<Button
		outline={chosen != 'Thursday'}
		onclick={() => {
			chosen = 'Thursday';
		}}
		class={`row-start-1 row-end-2 col-start-4 rounded-b-none rounded-r-none rounded-l-none ${chosen == 'Thursday' ? 'bg-accent' : ''}`}
		>Thursday</Button
	>
	<Button
		outline={chosen != 'Friday'}
		onclick={() => {
			chosen = 'Friday';
		}}
		class={`row-start-1 row-end-2 col-start-5 rounded-b-none rounded-r-none rounded-l-none ${chosen == 'Friday' ? 'bg-accent' : ''}`}
		>Friday</Button
	>
	<Button
		outline={chosen != 'Saturday'}
		onclick={() => {
			chosen = 'Saturday';
		}}
		class={`row-start-1 row-end-2 col-start-6 rounded-b-none rounded-r-none rounded-l-none ${chosen == 'Saturday' ? 'bg-accent' : ''}`}
		>Saturday</Button
	>
	<Button
		outline={chosen != 'Sunday'}
		onclick={() => {
			chosen = 'Sunday';
		}}
		class={`row-start-1 row-end-2 col-start-7 rounded-b-none rounded-r-2xl  rounded-l-none ${chosen == 'Sunday' ? 'bg-accent' : ''}`}
		>Sunday</Button
	>

	<div class="col-start-1 col-end-8 row-start-2 rounded-b-2xl row-end-8 bg-dark-compliment">
		{weekDaysDates[chosen]}
	</div>
</div>
