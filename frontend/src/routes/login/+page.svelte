<script lang="ts">
	import { applyAction, enhance } from '$app/forms';
	import '$lib/assets/styles/login.sass';
	import { setAuthToken } from '$lib/requestUtils/setAuthToken';

	// import type { PageProps } from './$types';
	// import { getContext } from 'svelte';

	//const sendBread: (e: unknown) => void = getContext('sendBread');

	///** @type {import('./$types').PageProps} */
	//let { form }: PageProps = $props();
	let errorMessage = $state('');
</script>

<div class="login-header">LOGIN</div>
<form
	class="login-body"
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
	<input id="email" type="email" name="email" placeholder="unbelivableemail@zmail.com" required />
	<input id="password" type="password" name="password" placeholder="123987" required />

	<p class="error">{errorMessage}</p>

	<button name="submit" type="submit"> SUBMIT </button>
</form>
