<script lang="ts">
	import { Card } from 'flowbite-svelte';
	import DOMpurify from 'dompurify';

	const { imageLink = '', caption, text } = $props();
</script>

<!-- /*got the joke? News are not countable*/ -->
<Card img={imageLink} class="new mb-2 break-inside-avoid min-w-0 w-full">
	<div class="m-6">
		<h5
			class="break-all caption mb-2 text-2xl font-bold bg-dark text-light font-hollow min-h-6 max-h-20 overflow-hidden"
		>
			{caption}
		</h5>
		<p
			class="break-all text mb-3 leading-tight font-normal bg-dark text-light font-hollow min-h-8 max-h-80 overflow-y-scroll no-scrollbar wrap-break-words mask-[linear-gradient(to_bottom,black_0%,black_80%,transparent_100%)]"
		>
			// eslint-disable-next-line svelte/no-at-html-tags
			{@html DOMpurify.sanitize(text)}
		</p>
	</div>
</Card>
