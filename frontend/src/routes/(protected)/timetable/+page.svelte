<script lang="ts">
	import {
		GetWeekDateFromWeek,
		GetSundayFromMonday,
		FormattedStringFromDate
	} from '$lib/timeUtils';
	import { ArrowRightOutline, ArrowLeftOutline } from 'flowbite-svelte-icons';

	import { Button, Input } from 'flowbite-svelte';
	import type { PageData } from './$types';
	import { getUserContext } from '$lib/context/userContext';
	import { enhance } from '$app/forms';
	import type { LessonCooked } from '$lib/types/timetable';

	let self = getUserContext();

	const { data }: { data: PageData } = $props();

	const dayNames = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday'];

	let weekNumber = $state((() => data.currentWeek)()); //change to onload week
	let currentYear = $state(2026);
	let currentDay = $state((() => data.currentDay)());

	let chosen = $state(currentDay); //change to server time + offset
	let weekDate = $derived.by((): string => {
		console.log('curryear+weekNumber: ', currentYear, weekNumber);
		const dateMonday = GetWeekDateFromWeek(currentYear, weekNumber);
		const formattedDateMonday = FormattedStringFromDate(dateMonday);

		const dateSunday = GetSundayFromMonday(dateMonday);
		const formattedDateSunday = FormattedStringFromDate(dateSunday);

		return `${formattedDateMonday}->${formattedDateSunday}`;
	});
	let formElement: HTMLFormElement;

	let lessons = $state(data.lessons);
	$inspect(lessons);
</script>

<div
	class="col-start-2 row-start-1 col-end-12 row-end-3 bg-dark text-[min(clamp(1rem,10vh,7rem),clamp(1rem,10vw,7rem))] inline-grid items-center justify-center italic font-hollow tracking-[0.5rem] select-none pointer-events-none"
>
	TIMETABLE
</div>

<div
	class="col-start-5 col-end-9 row-start-3 row-end-3 bg-dark -mb-8 rounded-t-4xl
	flex items-center justify-center
	"
>
	<Input value={weekNumber} name="weekNumber" type="hidden" />

	<Button pill={true} onclick={() => weekNumber--} class="p-2!" aria-label="turn to previous week"
		><ArrowLeftOutline class="h-6 w-6" /></Button
	>
	{weekDate}
	<Button pill={true} onclick={() => weekNumber++} class="p-2!" aria-label="turn to next week"
		><ArrowRightOutline class="h-6 w-6" /></Button
	>
</div>

<div
	class="
    col-start-1 row-start-4 row-end-12 col-end-13 bg-dark
    grid grid-cols-7 grid-rows-7 h-full w-full p-gutter
    rounded-4xl
    "
>
	{#each dayNames as day, idx (day)}
		<Button
			outline={currentDay != idx + 1}
			onclick={() => {
				currentDay = idx + 1;
				/* 				formElement.requestSubmit();
				 */
			}}
			class={[
				`row-start-1 row-end-2 col-start-${idx + 1} rounded-b-none`,
				currentDay === idx + 1 ? 'bg-accent' : '',
				idx === 6 ? 'rounded-r-2xl' : 'rounded-r-none',
				idx === 0 ? 'rounded-l-2xl' : 'rounded-l-none'
			]
				.filter(Boolean)
				.join(' ')}>{day}</Button
		>
	{/each}

	<div class="col-start-1 col-end-8 row-start-2 rounded-b-2xl row-end-8 bg-dark-compliment">
		{#if lessons}{:else}{/if}
	</div>
</div>

<form
	action="?/getTimetable"
	method="POST"
	//"Actions always use POST requests, since GET requests should never have side-effects."
	class="hidden"
	bind:this={formElement}
	use:enhance={() => {
		return async ({ update, result }) => {
			await update();
			console.log('res:', result);
			if (result.type === 'success' && result.data?.lessons) {
				lessons = result.data.lessons as LessonCooked[];
			}
		};
	}}
>
	<Input value={chosen} name="dayNumber" />
	<Input value={weekNumber} name="weekNumber" />
	<Input value={self.uuid} name="userUUID" />
</form>
