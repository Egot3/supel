<script lang="ts">
	//import '$lib/assets/global.scss';
	import '$lib/assets/styles/toaster.sass';
	import '../app.css';

	import { setContext } from 'svelte';
	import { SvelteMap } from 'svelte/reactivity';

	import Toaster from '$lib/components/Toaster.svelte';

	let { children } = $props();
	// svelte-ignore non_reactive_update
	//SvelteMap is reactive by default
	let messages = new SvelteMap<string, unknown>(); //uhm, ActUallY

	const sendBread = (message: unknown) => {
		console.log('adding element to map');
		messages.set(crypto.randomUUID(), message);
		console.log(messages);
	};

	setContext('sendBread', sendBread);
</script>

<div class="toaster">
	<Toaster bind:messages />
</div>
<div class="grid grid-cols-12 grid-rows-12 h-screen w-screen">
	<div class="avatar-wrapper"></div>
	{@render children()}
</div>

<style lang="sass">

</style>
