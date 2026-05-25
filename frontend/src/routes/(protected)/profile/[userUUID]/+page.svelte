<script lang="ts">
	import ProfileUpdateModal from '$lib/modals/ProfileUpdateModal.svelte';
	import { Button } from 'flowbite-svelte';
	import type { newCooked } from '$lib/types/new';
	import { innerWidth } from 'svelte/reactivity/window';
	import { enhance } from '$app/forms';
	import type { PageData } from './$types';
	import NewModal from '$lib/modals/NewModal.svelte';
	import NewExpanded from '$lib/components/NewExpanded.svelte';
	import AreYouSure from '$lib/modals/AreYouSure.svelte';
	import { type User } from '$lib/types/user';
	import { getUserContext } from '$lib/context/userContext';
	import UserCard from '$lib/components/UserCard.svelte';

	let { data }: { data: PageData } = $props();

	let page = $state(0);
	const size = 5;
	/* 	const { news, user } = $derived(data);
	 */
	const news = $derived(data.news);
	const user = $derived(data.user);
	let newsCooked: Array<newCooked> = $state(
		(() => {
			return news;
		})()
	);

	let updateProfileModalOpen = $state(false);

	let self = getUserContext();
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

	let isDeleting = $state(false);
</script>

<ProfileUpdateModal bind:open={updateProfileModalOpen} userInfo={(user ?? {}) as User} />
<NewModal bind:scoped bind:scopedNew />
<AreYouSure
	bind:open={isDeleting}
	actionName="account deletion"
	isReversable
	action="?/deleteUser"
	method="post"
	confirmationBody="enter your username"
	confirmationInput={user.nickname}
/>

<UserCard bind:isDeleting bind:updateProfileModalOpen {user} {self} />

<div
	class="overflow-scroll bg-forest-900 col-start-1 col-end-13 row-start-7 row-end-12 rounded-4x1 p-gutter gap-gutter"
	onscroll={UpdateRequest}
>
	{#each newsCooked as cooked (cooked.newId)}
		<NewExpanded
			newData={cooked}
			onclick={() => {
				scopedNew = cooked;
			}}
		/>
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
