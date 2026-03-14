<script lang="ts">
	//import '$lib/assets/styles/news.sass';
	import New from '$lib/components/New.svelte';
	import { idGenerator } from 'flowbite-svelte';
	const itemsRaw = [
		{ id: '1', text: '1', caption: '2' },
		{ id: '2', text: '1', caption: '2' },
		{ id: '3', text: '1', caption: '2' },
		{ id: '4', text: '1', caption: '2' },
		{
			id: '5',
			text: 'Cillum excepteur nisi exercitation do eu. Reprehenderit cupidatat ea consectetur ad tempor. Veniam non duis cupidatat in nulla esse occaecat enim officia nisi eu ullamco ipsum. Esse fugiat aliqua non sint. Deserunt est cillum ipsum adipisicing est ullamco.',
			caption: '2'
		},
		{ id: '6', text: '1', caption: '2' },
		{ id: '7', text: '1', caption: '2' },
		{ id: '8', text: '1', caption: '2' },
		{ id: '9', text: '1', caption: '2' },
		{ id: '10', text: '1', caption: '2' },
		{ id: '11', text: '1', caption: '2' },
		{ id: '12', text: '1', caption: '2' },
		{
			id: '13',
			text: 'Cillum excepteur nisi exercitation do eu. Reprehenderit cupidatat ea consectetur ad tempor. Veniam non duis cupidatat in nulla esse occaecat enim officia nisi eu ullamco ipsum. Esse fugiat aliqua non sint. Deserunt est cillum ipsum adipisicing est ullamco.',
			caption: '2'
		},
		{ id: '14', text: '1', caption: '2' },
		{ id: '15', text: '1', caption: '2' },
		{ id: '16', text: '1', caption: '2' },
		{ id: '17', text: '1', caption: '2' },
		{ id: '18', text: '1', caption: '2' },
		{ id: '19', text: '1', caption: '2' }
	];

	const dismantle = <T,>(inputArr: Array<T>, amount: number): Array<Array<T>> => {
		const arrOfDismantled: Array<Array<T>> = Array.from({ length: amount }, () => []);
		inputArr.forEach((element, idx) => {
			arrOfDismantled[idx % amount].push(element);
		});
		return arrOfDismantled;
	};

	const items = $derived(dismantle(itemsRaw, 4));
</script>

/* eslint-disable svelte/require-each-key */ /* eslint-disable svelte/require-each-key */
<div
	class="inline-flex bg-accent items-center justify-center text-5xl font-hollow italic cursor-pointer select-none row-start-1"
>
	POST
</div>

<div
	class="col-start-3 row-start-1 col-end-11 row-end-3 bg-dark text-[clamp(1rem,10vh,7rem)] inline-grid items-center justify-center italic font-hollow tracking-[0.5rem] select-none pointer-events-none"
>
	NEWS
</div>

<div class="bg-dark col-start-2 col-end-12 row-start-4 row-end-12 overflow-scroll">
	<!-- <Gallery class="grid-cols-2 md:grid-cols-4 gap-4 p-4">
		<Gallery>
			{#each items[0] as item (item.id)}
				<New caption={item.caption} text={item.text} />
			{/each}
		</Gallery>
		<Gallery>
			{#each items[1] as item (item.id)}
				<New caption={item.caption} text={item.text} />
			{/each}
		</Gallery>
		<Gallery>
			{#each items[2] as item (item.id)}
				<New caption={item.caption} text={item.text} />
			{/each}
		</Gallery>
		<Gallery>
			{#each items[3] as item (item.id)}
				<New caption={item.caption} text={item.text} />
			{/each}
		</Gallery>
	</Gallery> -->

	<div class="flex flex-row gap-4 p-4">
		{#each items as column}
			<div class="flex flex-col gap-4 flex-1">
				{#each column as item (item.id)}
					<New caption={item.caption} text={item.text} />
				{/each}
			</div>
		{/each}
	</div>
</div>
