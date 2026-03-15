<script lang="ts">
	import { applyAction, enhance } from '$app/forms';
	//import '$lib/assets/styles/login.sass';
	import { setAuthToken } from '$lib/requestUtils/setAuthToken';

	// import type { PageProps } from './$types';
	// import { getContext } from 'svelte';

	//const sendBread: (e: unknown) => void = getContext('sendBread');

	///** @type {import('./$types').PageProps} */
	//let { form }: PageProps = $props();
	let errorMessage = $state('');
</script>

<div
	class="min-h-0 min-w-0 bg-dark row-start-3 row-end-5 col-start-4 col-end-10 text-[clamp(1rem,10vh,7rem)] inline-grid items-center justify-center italic font-hollow tracking-[0.5rem] select-none pointer-events-none"
>
	LOGIN
</div>
<form
	class="bg-dark rounded-2x1 col-start-4 row-start-7 col-end-10 row-end-9 grid grid-cols-12 grid-rows-[repeat(2,1fr)_5%_10%] p-gutter gap-gutter rounded-2xl"
	action="?/login"
	method="POST"
	use:enhance={() => {
		//ah, yes (({})=>{({})=>{{}}})
		return async ({ result }) => {
			await applyAction(result);
			if (result.type === 'failure') {
				console.log(result.data);
				const e = result.data?.['errorMessage'];
				errorMessage = e ? String(e) : 'login failed';
			} else {
				if (result.type === 'success') {
					setAuthToken(String(result.data!['token']));
				}
			}
		};
	}}
>
	<input
		id="email"
		class="rounded-lg col-span-12 bg-dark-compliment font-hollow text-[0.8rem] text-light pl-2"
		type="email"
		name="email"
		placeholder="unbelivableemail@zmail.com"
		required
	/>
	<input
		id="password"
		class="rounded-lg col-span-12 bg-dark-compliment font-hollow text-[0.8rem] text-light pl-2"
		type="password"
		name="password"
		placeholder="123987"
		required
	/>

	<p
		id="error"
		class="flex items-center col-span-12 bg-dark font-hollow text-[0.8rem] text-accent pl-2"
	>
		{errorMessage}
	</p>

	<button
		name="submit"
		class="bg-accent rounded-lg font-bold text-dark font-hollow text-[0.8rem] col-start-5 col-end-9 items-center leading-0"
		type="submit"
	>
		SUBMIT
	</button>
</form>
