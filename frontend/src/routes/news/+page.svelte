<script lang="ts">
	//import '$lib/assets/styles/news.sass';
	import New from '$lib/components/New.svelte';
	import { withoutAuth } from '$lib/requestUtils/axiosConfigs';
	import axios from 'axios';
	import { Button, Spinner } from 'flowbite-svelte';

	interface newInterface {
		id: string;
		text: string;
		caption: string;
		imageLink?: string;
	}

	let page = 0;
	const size = 50;

	const itemsRaw: Array<newInterface> = $state([]);

	async function fetchNextPage() {
		loading = true;
		const cfg = withoutAuth(`http://localhost:5004/api/post?page=${page}&size=${size}`, 'get');
		axios(cfg)
			.then((res) => {
				itemsRaw.push(...res.data.items);
				console.log(res.data);
				if ((page + 1) * size < res.data.total) {
					page++;
				}
			})
			.catch((reason) => console.log(reason))
			.finally(() => (loading = false));
	}

	function dismantle(toDismantle: newInterface[], k: number): newInterface[][] {
		if (toDismantle.length == 0) {
			return [];
		}
		const sorted = [...toDismantle].sort((a, b) => b.text.length - a.text.length);
		const bins: newInterface[][] = Array.from({ length: k }, () => []);
		const sums = new Array(k).fill(0);

		for (const n of sorted) {
			let minIndex = 0;
			for (let i = 1; i < k; i++) {
				if (sums[i] < sums[minIndex]) minIndex = i;
			}
			bins[minIndex].push(n);
			let captionHeigth = (n.caption.length * 1.5) / 6;
			let textHeight = n.text.length / 6; //не спрашивайте, вообще без понятия

			if (captionHeigth > 20) {
				captionHeigth = 20;
			} else {
				if (captionHeigth < 6) {
					captionHeigth = 6;
				}
			}

			if (textHeight > 70) {
				textHeight = 70;
			} else {
				if (textHeight < 8) {
					textHeight = 8;
				}
			}

			sums[minIndex] += captionHeigth + textHeight + 1.5;
		}
		return bins;
	}

	let isAtBottom = false;
	let loading = $state(false);

	function handleScroll(event: UIEvent) {
		const element = event.currentTarget;
		const { scrollTop, scrollHeight, clientHeight } = element; // оно есть
		const atBottom = scrollTop + clientHeight >= scrollHeight - 10;

		if (atBottom !== isAtBottom) {
			isAtBottom = atBottom;
			if (atBottom) {
				fetchNextPage();
			}
		}
	}

	fetchNextPage();

	const items = $derived(dismantle(itemsRaw, 4));
</script>

<!-- <PostAction>
	POST
</PostAction> -->

<div
	class="col-start-3 row-start-1 col-end-11 row-end-3 bg-dark text-[clamp(1rem,10vh,7rem)] inline-grid items-center justify-center italic font-hollow tracking-[0.5rem] select-none pointer-events-none"
>
	NEWS
</div>

<div
	class="bg-dark col-start-2 col-end-12 row-start-4 row-end-12 overflow-scroll min-h-full inline-grid columns-1 grid-rows-1"
	// та самая сетка 1 на 1
	onscroll={handleScroll}
>
	{#if items.length}
		<div class="flex flex-row gap-gutter p-gutter min-w-0">
			{#each items as column, i (i)}
				<div class="flex flex-col gap-4 flex-1 min-w-0">
					{#each column as item (item.id)}
						<New caption={item.caption} text={item.text} />
					{/each}
				</div>
			{/each}
		</div>
	{:else}
		<div class="flex justify-center items-center self-center justify-self-center flex-col">
			no news
			<Button class="bg-accent p-2" size="xl" onclick={fetchNextPage}>
				{#if loading}
					<Spinner class="me-3" size="16" color="gray" />Loading ...
				{:else}
					Reload
				{/if}
			</Button>
		</div>
	{/if}
</div>
