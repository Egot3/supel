<script lang="ts">
	import { applyAction, enhance } from '$app/forms';
	import { setUserContext } from '$lib/context/userContext';
	import { Tooltip } from 'flowbite-svelte';

	let errorMessage = $state('');
	let errorType = $state('');
</script>

<div
	class="break-all min-h-0 min-w-0 bg-dark row-start-3 row-end-5 col-start-4 col-end-10 text-[min(clamp(1rem,10vh,7rem),clamp(1rem,10vw,7rem))] inline-grid items-center justify-center italic font-hollow tracking-[0.5rem] select-none pointer-events-none"
>
	LOGIN
</div>
<form
	class="overflow-scroll bg-dark rounded-2x1 col-start-4 row-start-7 col-end-10 row-end-9 grid grid-cols-12 grid-rows-[repeat(2,1fr)_15%] p-gutter gap-gutter rounded-2xl"
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
				const et = result.data?.['error'];
				errorType = String(et).includes('EMAIL') ? 'EMAIL' : 'PASSWORD';
				console.log(errorType);
			}
		};
	}}
>
	<div class="col-span-12">
		<label for="email">Email</label>
		<Tooltip class="bg-accent emailTooltip" isOpen={errorType === 'EMAIL'} reference="#email"
			>{errorMessage}</Tooltip
		>
		<input
			id="email"
			class="rounded-lg w-full h-[50%] bg-dark-compliment font-hollow text-[0.8rem] text-light pl-2"
			type="email"
			name="email"
			placeholder="unbelivableemail@zmail.com"
			required
		/>
	</div>

	<div class="col-span-12">
		<label for="password">Password</label>
		<Tooltip
			class="bg-accent passwordTooltip"
			isOpen={errorType === 'PASSWORD'}
			reference="#password">{errorMessage}</Tooltip
		>
		<input
			id="password"
			class="rounded-lg w-full h-[50%] bg-dark-compliment font-hollow text-[0.8rem] text-light pl-2"
			type="password"
			name="password"
			placeholder="123987"
			required
		/>
	</div>

	<!-- <p
		id="error"
		class="flex items-center col-span-12 bg-dark font-hollow text-[0.8rem] text-accent pl-2"
	>
		{errorMessage}
	</p> -->

	<button
		name="submit"
		class="bg-accent rounded-lg font-bold text-dark h-full font-hollow text-[0.8rem] col-start-5 col-end-9 items-center leading-0"
		type="submit"
	>
		SUBMIT
	</button>
</form>
