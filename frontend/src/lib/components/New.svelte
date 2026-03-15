<script lang="ts">
	import { Card } from 'flowbite-svelte';
	const { imageLink = '', caption, text } = $props();
</script>

<!-- /*got the joke? News are not countable*/ -->
<Card img={imageLink} class="new mb-2 break-inside-avoid max-w-[25%]">
	<div class="m-6">
		<h5 class="caption mb-2 text-2xl font-bold bg-dark text-light font-hollow min-h-6 max-h-20">
			{caption}
		</h5>
		<p
			class="wrap-break-word text mb-3 leading-tight font-normal bg-dark text-light font-hollow min-h-8 max-h-80 overflow-scroll no-scrollbar wrap-break-words mask-[linear-gradient(to_bottom,black_0%,black_80%,transparent_100%)]"
		>
			{text}
		</p>
	</div>
</Card>
