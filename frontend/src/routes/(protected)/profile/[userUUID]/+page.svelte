<script lang="ts">
	import { getUserContext } from '$lib/context/userContext';
	import ProfileUpdateModal from '$lib/modals/ProfileUpdateModal.svelte';
	import { Avatar, Button, P } from 'flowbite-svelte';
	import type { newCooked } from '$lib/types/new';
	import { innerWidth } from 'svelte/reactivity/window';
	import { enhance } from '$app/forms';
	import type { PageData } from './$types';
	import NewModal from '$lib/modals/NewModal.svelte';
	import NewExpanded from '$lib/components/NewExpanded.svelte';

	let { data }: { data: PageData } = $props();

	let page = $state(0);
	const size = 50;
	// svelte-ignore state_referenced_locally
	const startingNews = (data as { news: newCooked[] }).news;

	let newsCooked: Array<newCooked> = $state(startingNews);

	$inspect(newsCooked);

	let user = getUserContext();
	let updateProfileModalOpen = $state(false);

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

	let scopedNew = $state({} as newCooked);
	let scoped = $derived(Boolean(Object.keys(scopedNew).length > 0));
	$effect(() => {
		if (!scoped && Object.keys(scopedNew).length > 0) {
			scopedNew = {} as newCooked;
		}
	});
</script>

<ProfileUpdateModal bind:open={updateProfileModalOpen} userInfo={user} />
<NewModal bind:scoped bind:scopedNew />

<div
	class="bg-forest-900
    col-start-1 col-end-13 row-start-2 row-end-6
    flex
	p-gutter
	rounded-4xl
    "
>
	<div class="h-full w-auto mr-gutter">
		<Avatar class="h-full w-auto rounded-2xl" src={user.avatarUrl}></Avatar>
	</div>

	<div class="flex flex-col">
		<div class="h-full">
			<h1
				class="wrap-break-word caption mb-2 text-4xl font-bold bg-dark text-light font-hollow min-h-6 max-h-10 overflow-hidden"
			>
				{user.nickname}
			</h1>
			<P>{user.description}</P>
		</div>
		<div class="">
			<Button
				onclick={() => {
					updateProfileModalOpen = true;
				}}
				class="text-1x1 h-fit  self-end bg-accent">Change</Button
			>
		</div>
	</div>
</div>

<div
	class="overflow-scroll bg-forest-900 col-start-1 col-end-13 row-start-7 row-end-12 rounded-4x1 p-gutter gap-gutter"
	onscroll={UpdateRequest}
>
	{#each newsCooked as cooked (cooked.newId)}
		<NewExpanded newData={cooked} />
	{/each}
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
