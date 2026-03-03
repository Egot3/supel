<script lang="ts">
	import '$lib/assets/global.scss';
	import Toaster from '$lib/components/Toaster.svelte';
	import '$lib/assets/styles/toaster.sass';
	import { setContext } from 'svelte';
	let { children } = $props();
	let messages: Map<string, unknown> = $state(new Map<string, unknown>());
	const sendBread = (message: unknown) => {
		console.log('adding element to map');
		messages.set(crypto.randomUUID(), message);
		messages = new Map(messages);
		console.log(messages);
	};
	setContext('sendBread', sendBread);
</script>

<div class="toaster">
	<Toaster bind:messages />
</div>
<div
	style="display: grid; grid-template: repeat(12, 1fr) / repeat(12, 1fr); height: 100vh; width:100vw"
>
	{@render children()}
</div>
