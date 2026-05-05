<script lang="ts">
	import Moveable from 'svelte-moveable';
	import ImageFixerItem from './ImageFixerItem.svelte';

	let {
		width = $bindable(0),
		size = $bindable(0),
		height = $bindable(0),
		imageFiles = $bindable()
	}: {
		width: number;
		size: number;
		height: number;
		imageFiles: FileList | null;
	} = $props();
	let container: HTMLElement | null = $state(null);

	// eslint-disable-next-line svelte/prefer-svelte-reactivity
	let keyMap = new Map<File, string>();

	let imagesWithKeys = $derived(
		Array.from(imageFiles ?? []).map((file) => {
			if (!keyMap.has(file)) {
				keyMap.set(file, window.crypto.randomUUID());
			}
			return { file, key: keyMap.get(file)! };
		})
	);
</script>

<div class="container w-62.5 h-62.5" bind:this={container}>
	{#each imagesWithKeys as { file, key }, idx (key)}
		{console.log(file)}
		{#if idx == 0}
			<img src={URL.createObjectURL(file) ?? ''} class="" alt="" />
		{:else}
			<ImageFixerItem {file} />
		{/if}
	{/each}
</div>
{#if container}
	<Moveable
		target={container}
		clipTargetBounds={true}
		clippable={true}
		keepRatio
		defaultClipPath="circle"
		on:clip={({ detail: e }) => {
			e.target.style.clipPath = e.clipStyle;
			size = Number(e.clipStyles[0].slice(0, -2)); //Source?
			width = Number(e.clipStyles[2].slice(0, -2)); // trust me bro
			height = Number(e.clipStyles[3].slice(0, -2));

			console.log('width: %d, height: %d, size: %d, all:', width, height, size, e.clipStyle);
		}}
		on:drag={({ detail: e }) => {
			e.target.style.transform = e.transform;
		}}
	/>
{/if}
