<script lang="ts">
	import New from '$lib/components/New.svelte';
	import { type newCooked } from '$lib/types/new';
	import { Button } from 'flowbite-svelte';

	import { enhance } from '$app/forms';
	import { innerWidth } from 'svelte/reactivity/window';

	import type { PageData } from './$types';
	import NewModal from '$lib/modals/NewModal.svelte';
	import CreateNewModal from '$lib/modals/CreateNewModal.svelte';
	import { browser } from '$app/environment';

	let { data }: { data: PageData } = $props();

	let page = $state(0);
	const size = 50;

	let newsCooked: Array<newCooked> = $state(
		(() => {
			return data.news;
		})()
	);

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
		const { scrollTop } = element as HTMLElement;
		const { scrollHeight, clientHeight } = element as HTMLElement; // оно есть
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

	/* let user = getUserContext(); */
</script>

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
	{#if items.length && browser}
		<div class="flex flex-row gap-gutter p-gutter min-w-0">
			{#each items as column, i (i)}
				<div class="flex flex-col gap-4 flex-1 min-w-0">
					{#each column as item (item.newId)}
						<New
							onclick={() => {
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
