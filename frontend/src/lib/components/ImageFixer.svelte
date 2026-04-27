<script lang="ts">
	import Moveable from 'svelte-moveable';
	import ImageFixerItem from './ImageFixerItem.svelte';

	let {
		width = $bindable(0),
		height = $bindable(0),
		imageFiles = $bindable()
	}: {
		width: number;
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
		}}
		on:drag={({ detail: e }) => {
			e.target.style.transform = e.transform;
		}}
	/>
{/if}
