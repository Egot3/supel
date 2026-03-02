<script lang="ts">
	import '$lib/assets/global.scss';
	import Toaster from '$lib/components/Toaster.svelte';
	import '$lib/assets/styles/toaster.sass';
	import { setContext } from 'svelte';
	let { children } = $props();
	let messages: Array<unknown> = $state([]);
	const sendBread = (message: unknown) => {
		messages.push(message);
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
