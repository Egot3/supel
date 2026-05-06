<script lang="ts">
	import type { newCooked } from '$lib/types/new';
	import { Carousel, Controls, ControlButton, P, Modal } from 'flowbite-svelte';
	import { fly } from 'svelte/transition';
	import DOMPurify from 'isomorphic-dompurify';
	import type { User } from '$lib/types/user';
	import { Carouselify } from '$lib/compatUtils/carouselify';

	let { scoped = $bindable(false), scopedNew = $bindable({} as newCooked) } = $props();

	let userInfo: User = $state({} as User);
	async function GetUser() {
		const res = await fetch(`/api/user/${scopedNew.userId}`, {
			method: 'GET'
		});
		const data = await res.json();
		userInfo = data.user;
	}

	$effect(() => {
		if (scoped) {
			GetUser();
			/* 			console.log('user name:', userInfo.nickname);
			 */
		}
	});
	/* $inspect(userInfo.nickname); */
	const carouselItems = $derived(Carouselify(scopedNew.imageUrls ?? []));
</script>

<Modal
	bind:open={scoped}
	size="xs"
	class="bg-forest-900 dark:bg-forest-900 backdrop:bg-linen-900/50 text-linen-200
	p-gutter -mt-20 justify-self-center self-center
	overflow-scroll leading-[-1rem]
	h-[50dvh] min-h-100 
	"
	transition={fly}
	dismissable={false}
>
	<div class="grid grid-cols-12 min-h-full">
		<div class="col-start-1 col-end-9 grid grid-cols-3 gap-gutter">
			<div class="col-start-1 col-end-2 flex w-full h-full justify-center items-center">
				{#if scopedNew.imageUrls.length > 0}
					{#key carouselItems}
						{$inspect('carousel items: ', carouselItems)}
						{console.log('new has images')}
						<Carousel
							images={carouselItems}
							classes={{ slide: 'w-fit h-fit ' }}
							class="relative w-full h-70 flex justify-center items-center"
						>
							<Controls>
								{#snippet children(changeSlide)}
									<ControlButton
										name="Previous"
										forward={false}
										onclick={() => changeSlide(false)}
									/>
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
					{scopedNew.caption}
				</h5>
				<h6
					class="wrap-break-word caption mb-2 text-1xl font-bold bg-dark text-light font-hollow min-h-6 max-h-20 overflow-hidden"
				>
					{userInfo.nickname}
				</h6>

				<P class="bg-forest-900 dark:bg-forest-900 text-linen-200 w-[50%]">
					<!-- eslint-disable-next-line svelte/no-at-html-tags -->
					{@html DOMPurify.sanitize(scopedNew.body ?? '')}</P
				>
			</span>
		</div>
		<div class="col-start-9 col-end-13">2</div>
	</div>
</Modal>
