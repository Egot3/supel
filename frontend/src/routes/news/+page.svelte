<script lang="ts">
	//import '$lib/assets/styles/news.sass';
	import New from '$lib/components/New.svelte';
	import { withoutAuth } from '$lib/requestUtils/axiosConfigs';
	import { redirect } from '@sveltejs/kit';

	import axios from 'axios';

	import {
		Button,
		Spinner,
		Modal,
		Label,
		Input,
		Helper,
		Textarea,
		Fileupload
	} from 'flowbite-svelte';
	import { MailBoxOutline } from 'flowbite-svelte-icons';
	import { slide } from 'svelte/transition';

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
		const cfg = withoutAuth(`http://localhost:5004/api/post?page=${page}&size=${size}`, 'get', {});
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

	let popupModal = $state(false);

	let captionText = $state('');
	let bodyText = $state('');
	// let imagePath = $state(''); TBU
</script>

<Button
	onclick={() => (popupModal = true)}
	name="postingButton"
	class="bg-coral-500 col-start-1 text-4xl italic">New new</Button
>
<!-- my genius knows no border-radius -->

<Modal
	bind:open={popupModal}
	size="xs"
	class="bg-forest-900 dark:bg-forest-900 backdrop:bg-linen-900/50 text-linen-200 p-gutter mt-12 justify-self-center"
	transition={slide}
	dismissable={false}
>
	<form action="?/post" method="POST">
		<h2 class="mb-5 text-lg font-normal flex text-light bg-dark">
			<MailBoxOutline class="mx-left mb-2 h-6 w-6 text-light bg-dark" />
			New constructor
		</h2>

		<div class="text-left space-y-0.5 grid grid-cols-[65%_35%]">
			<div>
				<Label for="caption">Caption</Label>
				<Input
					name="caption"
					placeholder="extremly cool and original caption"
					class="rounded-lg bg-dark-compliment font-hollow text-[0.8rem] text-light pl-2"
					bind:value={captionText}
				/>
				<Helper></Helper>

				<Label for="text">Body</Label>
				<Textarea
					name="text"
					placeholder="even cooler, originaler and longer text"
					class="rounded-lg w-full bg-dark-compliment font-hollow text-[0.8rem] text-light pl-2 resize-none"
					rows={5}
					bind:value={bodyText}
				/>
				<Helper></Helper>

				<Label for="image">Uploaded image</Label>
				<Fileupload
					type="image"
					name="image"
					class="rounded-lg bg-dark-compliment font-hollow text-[0.8rem] text-light pl-2"
				>
					<!-- <img /> -->
				</Fileupload>
			</div>
			<div class="ml-10 mr-10">
				<New
					caption={captionText ? captionText : 'Preview'}
					text={bodyText ? bodyText : 'Your text could be here'}
				/>
			</div>
		</div>

		<div class="space-x-2 bg-dark text-light mt-6">
			<Button
				type="submit"
				value="post"
				name="goPostIt"
				class="text-center font-medium inline-flex items-center justify-center text-linen-200 bg-coral-700 hover:bg-coral-800 dark:bg-coral-600 dark:hover:bg-coral-700 focus-within:ring-coral-300 dark:focus-within:ring-coral-900 px-5 py-2.5 text-sm focus-within:ring-4 focus-within:outline-hidden rounded-lg"
				color="red">Post now</Button
			>
			<Button
				// type="cancel"
				value="cancel"
				onclick={() => (popupModal = false)}
				class="text-center font-medium inline-flex items-center justify-center text-linen-900 bg-transparent border border-linen-200 dark:border-linen-600 hover:bg-linen-100 dark:bg-linen-800 dark:text-linen-400 hover:text-primary-700 focus-within:text-primary-700 dark:focus-within:text-linen-50 dark:hover:text-linen-50 dark:hover:bg-linen-700 focus-within:ring-linen-200 dark:focus-within:ring-linen-700 px-5 py-2.5 text-sm focus-within:ring-4 focus-within:outline-hidden rounded-lg"
				color="dark">cancel</Button
			>
		</div>
	</form>
</Modal>

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
			<Button class="bg-accent p-2 w-l" size="xl" onclick={fetchNextPage}>
				{#if loading}
					<Spinner class="me-3" size="4" color="gray" />Loading...
				{:else}
					Reload
				{/if}
			</Button>
		</div>
	{/if}
</div>
