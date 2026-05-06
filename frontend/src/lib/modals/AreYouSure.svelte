<script lang="ts">
	import { Button, Input, Modal, P } from 'flowbite-svelte';
	import { fly } from 'svelte/transition';

	let {
		open = $bindable(false),
		actionName,
		action,
		method = 'post' as 'post' | 'get',
		confirmationBody,
		confirmationInput,
		isReversable = false,
		isFavourable = false
	} = $props();

	let userConfimation = $state('');
	let isConfirmed = $derived(userConfimation === confirmationInput);
</script>

<Modal
	bind:open
	size="xs"
	class="bg-forest-900 dark:bg-forest-900 backdrop:bg-linen-900/50 text-linen-200
	p-gutter -mt-20 justify-self-center self-center
	overflow-scroll leading-[-1rem]
	h-[25dvh] min-h-50 
	"
	transition={fly}
	dismissable={false}
>
	<form {method} {action}>
		<h1>Stop right there!</h1>
		<P>
			The action you are about to commit({actionName}) is
			{isReversable ? 'reversable(but pretty hard to recover)' : 'not reversable'}.
		</P>
		<P>Are you sure you don't need a trusted adult to do it?</P>
		<P>{confirmationBody}</P>
		<Input placeholder={confirmationInput} bind:value={userConfimation} />

		{#if isFavourable}
			<Button type="submit" disabled={!isConfirmed} class="bg-accent">I am sure</Button>
			<Button
				onclick={() => {
					open = false;
				}}
				class="bg-dark-compliment"
			>
				Return
			</Button>
		{:else}
			<Button
				onclick={() => {
					open = false;
				}}
				class="bg-accent"
			>
				Return
			</Button>
			<Button type="submit" disabled={!isConfirmed} class="bg-dark-compliment">I am sure</Button>
		{/if}
	</form>
</Modal>
