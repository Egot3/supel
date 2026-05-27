<script lang="ts">
	import { onMount } from 'svelte';
	import '../app.css';
	let { children } = $props();

	const GFUEL = ['g', 'f', 'v', 'e', 'l', 'l', 'l', 'l'];

	let buffer = $state([] as string[]);
	let activated = $state(false);
	let coords = $state({ x: 0, y: 0 });

	onMount(() => {
		const handler = (e: KeyboardEvent) => {
			buffer = [...buffer.slice(-(GFUEL.length - 1)), e.key];
			if (buffer.join(',') === GFUEL.join(',')) {
				activated = true;
				setTimeout(() => (activated = false), 1200);
			}
		};

		window.addEventListener('keydown', handler);

		function handlePointerMove(e: PointerEvent) {
			coords.x = e.clientX;
			coords.y = e.clientY;
		}
		window.addEventListener('pointermove', handlePointerMove);

		return () => {
			window.removeEventListener('keydown', handler);
			window.removeEventListener('pointermove', handlePointerMove);
		};
	});
	//$inspect(coords);
</script>

<div class="grid grid-cols-12 grid-rows-12 gap-gutter">
	{@render children()}
	{#if activated}
		<img
			class="absolute"
			alt="you caught an easter egg"
			src="/gifs/explosion-gif-transparent.gif"
			style="top: {coords.y - 100}px; left: {coords.x - 50}px;"
		/>
	{/if}
</div>
