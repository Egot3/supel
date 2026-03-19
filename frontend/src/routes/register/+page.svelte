<script lang="ts">
	import { applyAction, enhance } from '$app/forms';
	//import '$lib/assets/styles/register.sass';
	import '$lib/requestUtils/setAuthToken';

	let errorMessage = $state('');
</script>

<div
	class="min-h-0 min-w-0 bg-dark row-start-3 row-end-5 col-start-4 col-end-10 text-[clamp(1rem,10vh,7rem)] inline-grid items-center justify-center italic font-hollow tracking-[0.5rem] select-none pointer-events-none"
>
	REGISTER
</div>
<form
	class="bg-dark rounded-2x1 col-start-4 row-start-7 col-end-10 row-end-10 grid grid-cols-12 grid-rows-[repeat(4,1fr)_5%_10%] p-gutter gap-gutter rounded-2xl"
	action="?/register"
	method="POST"
	use:enhance={() => {
		//ah, yes (({})=>{({})=>{{}}})
		return async ({ result }) => {
			console.log('applying action');
			await applyAction(result);
			console.log('applied');
			if (result.type === 'failure') {
				console.log(result.data);
				const e = result.data?.['errorMessage'];
				errorMessage = e ? String(e) : 'login failed';
			} else {
				console.log('still hope');
			}
		};
	}}
>
	<input
		id="nickname"
		class="rounded-lg col-span-12 bg-dark-compliment font-hollow text-[0.8rem] text-light pl-2"
		name="nickname"
		placeholder="Your name"
	/>
	<input
		id="email"
		type="email"
		class="rounded-lg col-span-12 bg-dark-compliment font-hollow text-[0.8rem] text-light pl-2"
		name="email"
		placeholder="unbelivableemail@zmail.com"
		required
	/>
	<input
		id="password"
		type="password"
		class="rounded-lg col-span-12 bg-dark-compliment font-hollow text-[0.8rem] text-light pl-2"
		name="password"
		placeholder="123987"
		required
	/>
	<input
		class="rounded-lg col-span-12 bg-dark-compliment font-hollow text-[0.8rem] text-light pl-2"
		id="passwordDup"
		type="password"
		name="passwordDup"
		placeholder="123987"
		required
	/>

	<p
		class="flex rounded-lg items-center col-span-12 bg-dark font-hollow text-[0.8rem] text-accent pl-2"
		id="error"
	>
		{errorMessage}
	</p>

	<button
		type="submit"
		name="submit"
		class="bg-accent rounded-lg font-bold text-dark font-hollow text-[0.8rem] col-start-5 col-end-9 items-center leading-0"
	>
		SUBMIT
	</button>
</form>
