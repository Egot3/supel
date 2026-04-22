<script lang="ts">
	//import '$lib/assets/styles/news.sass';
	import New from '$lib/components/New.svelte';
	import { type newCooked } from '$lib/types/new';
	import {
		Button,
		// Spinner,
		Modal,
		Label,
		Input,
		Helper,
		Textarea,
		Fileupload,
		Carousel,
		Controls,
		P
	} from 'flowbite-svelte';
	import DOMPurify from 'isomorphic-dompurify';

	import { MailBoxOutline } from 'flowbite-svelte-icons';
	import { fly, slide } from 'svelte/transition';
	import { applyAction, enhance } from '$app/forms';
	import { innerWidth } from 'svelte/reactivity/window';
	import NewExample from '$lib/components/NewExample.svelte';
	import type { PageData } from '../$types';

	let { data }: { data: PageData } = $props();

	let page = $state(0);
	const size = 50;
	// svelte-ignore state_referenced_locally
	const startingNews = (data as { news: newCooked[] }).news;

	let newsCooked: Array<newCooked> = $state(startingNews);
	// newsCooked.push(...data.news);
	// console.log('news cooked: ', newsCooked);

	function dismantle(toDismantle: newCooked[], k: number): newCooked[][] {
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

	let createNewModal = $state(false);

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
</script>

<Button
	onclick={() => (createNewModal = true)}
	name="postingButton"
	class="bg-coral-500 col-start-1 text-[clamp(1rem,3vh,3rem)] italic">New new</Button
>
<!-- my genius knows no border-radius -->

<!-- {console.log('data: ', data)} -->
<!-- {console.log('newsCooked: ', newsCooked)} -->

<Modal
	bind:open={scoped}
	size="xs"
	class="bg-forest-900 dark:bg-forest-900 backdrop:bg-linen-900/50 text-linen-200
	p-gutter -mt-20 justify-self-center self-center
	overflow-scroll leading-[-1rem]
	h-[50dvh] min-h-100 
	"
	transition={fly}
	dismissable={false}
>
	<div class="grid grid-cols-12">
		<div class="col-start-1 col-end-9">
			{#if scopedNew.imageUrls.length > 0}
				<Carousel
					images={scopedNew.imageUrls.map((url) => {
						return { src: url };
					})}
				>
					<Controls />
				</Carousel>
			{/if}

			<h5
				class="wrap-break-word caption mb-2 text-2xl font-bold bg-dark text-light font-hollow min-h-6 max-h-20 overflow-hidden"
			>
				{scopedNew.caption}
			</h5>

			<P class="bg-forest-900 dark:bg-forest-900 text-linen-200">
				<!-- eslint-disable-next-line svelte/no-at-html-tags -->
				{@html DOMPurify.sanitize(scopedNew.body ?? '')}</P
			>
		</div>
		<div class="col-start-9 col-end-13">2</div>
	</div>
</Modal>

<Modal
	bind:open={createNewModal}
	size="xs"
	class="bg-forest-900 dark:bg-forest-900 backdrop:bg-linen-900/50 text-linen-200 p-gutter mt-12 justify-self-center overflow-scroll leading-[-1rem]"
	transition={slide}
	dismissable={false}
>
	<form
		action="?/post"
		method="POST"
		enctype="multipart/form-data"
		use:enhance={() => {
			//ah, yes (({})=>{({})=>{{}}})
			return async ({ result }) => {
				await applyAction(result);
				if (result.type === 'failure') {
					console.log(result.data);
					const e = result.data?.['errorMessage'];
					console.log(e);
				} else {
					createNewModal = false;
					bodyText = '';
					captionText = '';
				}
			};
		}}
	>
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
					required
					name="textArea"
					placeholder="even cooler, originaler and longer text"
					class="rounded-lg w-full bg-dark-compliment font-hollow text-[0.8rem] text-light pl-2 resize-none"
					rows={5}
					bind:value={bodyText}
				/>
				<Helper></Helper>

				<!-- {/* <Label for="image">Uploaded image</Label>
				<Fileupload
					type="image"
					name="image"
				>
				</Fileupload>
				<Helper>SVG, PNG, JPG or GIF (MAX. 800x400px).</Helper> * /} -->

				<Label for="with_helper" class="pb-2">Upload file</Label>
				<Fileupload
					name="image"
					id="with_helper"
					class="rounded-lg bg-dark-compliment font-hollow text-[0.8rem] text-light pl-2"
				/>
				<Helper class="bg-forest-900 text-linen-200 text-[0.75rem]"
					>SVG, PNG, JPG or GIF (MAX. 800x400px).</Helper
				>
			</div>
			<div class="ml-10 mr-10">
				<NewExample
					caption={captionText ?? 'Preview'}
					text={bodyText ?? 'Your text could be here'}
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
				onclick={() => (createNewModal = false)}
				class="text-center font-medium inline-flex items-center justify-center text-linen-900 bg-transparent border border-linen-200 dark:border-linen-600 hover:bg-linen-100 dark:bg-linen-800 dark:text-linen-400 hover:text-primary-700 focus-within:text-primary-700 dark:focus-within:text-linen-50 dark:hover:text-linen-50 dark:hover:bg-linen-700 focus-within:ring-linen-200 dark:focus-within:ring-linen-700 px-5 py-2.5 text-sm focus-within:ring-4 focus-within:outline-hidden rounded-lg"
				color="dark">cancel</Button
			>
		</div>
	</form>
</Modal>

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
