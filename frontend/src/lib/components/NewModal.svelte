<script lang="ts">
	import type { newCooked } from '$lib/types/new';
	import { Carousel, CarouselIndicators, P, Modal } from 'flowbite-svelte';
	import { fly } from 'svelte/transition';
	import DOMPurify from 'isomorphic-dompurify';
	import type { User } from '$lib/types/user';

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
			<div class="col-start-1 col-end-2 grid justify-center align-center">
				{#if scopedNew.imageUrls.length > 0}
					<div class="relative self-center">
						<Carousel
							images={scopedNew.imageUrls.map((url) => {
								return { src: url };
							})}
						>
							<CarouselIndicators />
						</Carousel>
					</div>
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
