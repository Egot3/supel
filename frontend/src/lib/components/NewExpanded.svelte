<script lang="ts">
	import type { newCooked } from '$lib/types/new';
	import { P, Carousel, Controls, Button, ControlButton } from 'flowbite-svelte';
	import DOMPurify from 'isomorphic-dompurify';

	let { newData }: { newData: newCooked } = $props();
	function Carouselify(imageSrcs: string[]) {
		return imageSrcs.map((imgSrc) => {
			return {
				src: imgSrc,
				alt: 'you could see something there',
				title: 'yes'
			};
		});
	}
	const carouselItems = $derived(Carouselify(newData.imageUrls));
	console.log('carousel items: ', carouselItems);
</script>

<div class="grid grid-cols-12 h-75 bg-forest-950">
	<div class="col-start-1 col-end-9 grid grid-cols-3 gap-gutter">
		<div class="col-start-1 col-end-3 flex w-full h-full justify-center items-center">
			{console.log('newData:', newData.imageUrls.length)}
			{#if newData.imageUrls.length > 0}
				{#key carouselItems}
					{console.log('new has images')}
					<Carousel
						images={carouselItems}
						classes={{ slide: 'w-fit h-fit ' }}
						class="relative w-full h-70 flex justify-center items-center"
					>
						<Controls>
							{#snippet children(changeSlide)}
								<ControlButton name="Previous" forward={false} onclick={() => changeSlide(false)} />
								<ControlButton name="Next" forward={true} onclick={() => changeSlide(false)} />
							{/snippet}
						</Controls>
					</Carousel>
				{/key}
			{/if}
		</div>

		<span>
			<h5
				class="wrap-break-word caption mb-2 text-2xl font-bold bg-dark text-light font-hollow min-h-6 max-h-20 overflow-hidden"
			>
				{newData.caption}
			</h5>

			<P class="bg-forest-900 dark:bg-forest-900 text-linen-200 w-[50%]">
				<!-- eslint-disable-next-line svelte/no-at-html-tags -->
				{@html DOMPurify.sanitize(newData.body ?? '')}</P
			>
		</span>
	</div>
</div>
