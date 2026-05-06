<script lang="ts">
	import { Avatar, Button, P } from 'flowbite-svelte';

	let {
		user,
		self,
		updateProfileModalOpen = $bindable(false),
		isDeleting = $bindable(false)
	} = $props();

	$inspect('user in component:', user);
</script>

<div
	class="bg-forest-900
    col-start-1 col-end-13 row-start-2 row-end-6
    flex
	p-gutter
	rounded-4xl
    "
>
	<div class="h-full w-auto mr-gutter">
		{console.log('user avatar: ', user.avatarUrl)}
		<Avatar class="h-full w-auto rounded-2xl" src={user.avatarUrl} />
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

		{#if self.uuid === user.uuid}
			<div class="">
				<Button
					onclick={() => {
						updateProfileModalOpen = true;
					}}
					class="text-1x1 h-fit  self-end bg-accent">Change</Button
				>
				<Button
					onclick={() => {
						isDeleting = true;
					}}>Delete</Button
				>
			</div>
		{/if}
	</div>
</div>
