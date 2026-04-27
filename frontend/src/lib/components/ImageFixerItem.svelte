<script lang="ts">
	let { file }: { file: File } = $props();
	import Moveable from 'svelte-moveable';

	let el = $state<HTMLElement | undefined>(undefined);
</script>

<div class="container">
	<div bind:this={el} class="h-fit w-fit absolute">
		<img src={URL.createObjectURL(file) ?? ''} class="pointer-none select-none" alt="" />
	</div>

	{#if el}
		{console.log('el exists')}
		<Moveable
			target={el}
			draggable={true}
			on:drag={({ detail: e }) => {
				e.target.style.transform = e.transform;
			}}
			scalable={true}
			keepRatio={true}
			on:scale={({ detail: e }) => {
				e.target.style.transform = e.drag.transform;
			}}
			snappable
			snapContainer=".container"
		/>
	{/if}
</div>
