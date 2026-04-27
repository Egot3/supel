<script lang="ts">
	//import '$lib/assets/styles/news.sass';
	import New from '$lib/components/New.svelte';
	import { type newCooked } from '$lib/types/new';
	import {
		Button,
		// Spinner,
		Avatar
	} from 'flowbite-svelte';

	import { enhance } from '$app/forms';
	import { innerWidth } from 'svelte/reactivity/window';

	import type { PageData } from '../$types';
	import NewModal from '$lib/modals/NewModal.svelte';
	import { getUserContext } from '$lib/context/userContext';
	import CreateNewModal from '$lib/modals/CreateNewModal.svelte';

	let { data }: { data: PageData } = $props();

	let page = $state(0);
	const size = 50;
	// svelte-ignore state_referenced_locally
	const startingNews = (data as { news: newCooked[] }).news;

	let newsCooked: Array<newCooked> = $state(startingNews);
	// newsCooked.push(...data.news);
	// console.log('news cooked: ', newsCooked);

	/* function dismantle(toDismantle: newCooked[], k: number): newCooked[][] {
		if (toDismantle.length == 0) {
			return [];
		}
		const sorted = [...toDismantle].sort((a, b) => b.bodySize - a.bodySize);
		const bins: newCooked[][] = Array.from({ length: k }, () => []);
		console.log('length of an arr: ', k);
		const sums = new Array(k).fill(0);

		for (const n of sorted) {
			let minIndex = 0;
			for (let i = 1; i < k; i++) {
				if (sums[i] < sums[minIndex]) minIndex = i;
			}
			bins[minIndex].push(n);
			let captionHeigth = (n.caption.length * 1.5) / 6;
			let textHeight = n.bodySize / 6; //не спрашивайте, вообще без понятия

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
	} */ //слетели все бинды(

	function dismantle(toDismantle: newCooked[], k: number): newCooked[][] {
		const resp: newCooked[][] = Array.from({ length: k }, () => []);

		toDismantle.forEach((item, idx) => {
			const target = idx % k;
			resp[target].push(item);
		});

		return resp;
	}

	let isAtBottom = false;
	let loading = $state(false);

	let formElement: HTMLFormElement;

	function UpdateRequest(event: UIEvent) {
		console.log('update trigger');
		if (loading) return;
		console.log('entered update');
		const element = event.currentTarget;
		const { scrollTop, scrollHeight, clientHeight } = element as HTMLElement; // оно есть
		const atBottom = scrollTop + clientHeight >= scrollHeight - 10;

		if (atBottom !== isAtBottom) {
			isAtBottom = atBottom;
			if (atBottom) {
				formElement.requestSubmit();
			}
		}
	}

	console.log('innerwidth:', innerWidth.current);
	const items = $derived(dismantle(newsCooked, Math.ceil((innerWidth.current ?? 2000) / 400)));

	let createNewModalOpen = $state(false);

	let captionText = $state('');
	let bodyText = $state('');
	// let imagePath = $state(''); TBU

	let scopedNew = $state({} as newCooked);
	let scoped = $derived(Boolean(Object.keys(scopedNew).length > 0));
	$effect(() => {
		if (!scoped && Object.keys(scopedNew).length > 0) {
			scopedNew = {} as newCooked;
		}
	});

	let user = getUserContext();
</script>

<Avatar src={user.avatarUrl} class="flex items-center justify-center col-start-12 col-end-13"
></Avatar>

<Button
	onclick={() => (createNewModalOpen = true)}
	name="postingButton"
	class="bg-coral-500 col-start-1 text-[clamp(1rem,3vh,3rem)] italic">New new</Button
>
<!-- my genius knows no border-radius -->

<!-- {console.log('data: ', data)} -->
<!-- {console.log('newsCooked: ', newsCooked)} -->

<NewModal bind:scoped bind:scopedNew />
<CreateNewModal bind:open={createNewModalOpen} bind:bodyText bind:captionText />

<div
	class="col-start-3 row-start-1 col-end-11 row-end-3 bg-dark text-[min(clamp(1rem,10vh,7rem),clamp(1rem,10vw,7rem))] inline-grid items-center justify-center italic font-hollow tracking-[0.5rem] select-none pointer-events-none"
>
	NEWS
</div>

<div
	class="bg-dark col-start-2 col-end-12 row-start-4 row-end-12 overflow-scroll min-h-full inline-grid columns-1 grid-rows-1"
	// та самая сетка 1 на 1
	onscroll={UpdateRequest}
>
	{#if items.length}
		<div class="flex flex-row gap-gutter p-gutter min-w-0">
			{#each items as column, i (i)}
				<div class="flex flex-col gap-4 flex-1 min-w-0">
					{#each column as item (item.newId)}
						<New
							onclick={() => {
								//malpractice
								scopedNew = item;
							}}
							newData={item}
						/>
					{/each}
				</div>
			{/each}
		</div>
	{:else}
		<!-- <div class="flex justify-center items-center self-center justify-self-center flex-col">
			no news
			<Button class="bg-accent p-2 w-l" size="xl" name="reload" 
			//onclick={UpdateRequest}
			>
				{#if loading}
					<Spinner class="me-3" size="4" color="gray" />Loading...
				{:else}
					Reload
				{/if}
			</Button>
		</div> -->
	{/if}
	<form
		bind:this={formElement}
		method="POST"
		action="?/loadMore"
		use:enhance={() => {
			loading = true;
			console.log('entered loading', loading);

			return async ({ update, result }) => {
				await update();
				console.log('res:', result);
				if (result.type === 'success' && result.data?.news) {
					newsCooked = [...newsCooked, ...(result.data.news as newCooked[])];
					page++;
				}
				loading = false;
			};
		}}
	>
		<input type="hidden" name="page" value={page + 1} />
		<input type="hidden" name="size" value={size} />
		<Button type="submit" class="invisible">{loading ? 'deeper we go' : 'fetch more news'}</Button>
	</form>
</div>
