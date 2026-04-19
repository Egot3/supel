<script lang="ts">
	import type { newProps } from '$lib/types/new';
	import { Card } from 'flowbite-svelte';
	import DOMPurify from 'isomorphic-dompurify';
	const { newData, onclick }: newProps = $props();
</script>

<!-- /*got the joke? News are not countable*/ -->
<Card {onclick} img={newData.imageUrls[0] ?? ''} class="new mb-2 break-inside-avoid min-w-0 w-full">
	<div class="m-6">
		<h5
			class="break-all caption mb-2 text-2xl font-bold bg-dark text-light font-hollow min-h-6 max-h-20 overflow-hidden"
		>
			{newData.caption}
		</h5>
		{#if newData.body}
			<p
				class="break-all text mb-3 leading-tight font-normal bg-dark text-light font-hollow min-h-8 max-h-80 overflow-y-scroll no-scrollbar wrap-break-words mask-[linear-gradient(to_bottom,black_0%,black_80%,transparent_100%)]"
			>
				<!-- eslint-disable-next-line svelte/no-at-html-tags -->
				{@html DOMPurify.sanitize(newData.body)}
			</p>
		{/if}
	</div>
</Card>
