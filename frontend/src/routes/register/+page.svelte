<script lang="ts">
	import { applyAction, enhance } from '$app/forms';
	import '$lib/assets/styles/register.sass';
	import '$lib/requestUtils/setAuthToken';
	import { setAuthToken } from '$lib/requestUtils/setAuthToken';
	import { redirect } from '@sveltejs/kit';

	let errorMessage = $state('');
</script>

<div class="register-header">REGISTER</div>
<form
	class="register-body"
	action="?/register"
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
					redirect(305, '/');
				}
			}
		};
	}}
>
	<input id="nickname" name="nickname" placeholder="Your name" />
	<input id="email" type="email" name="email" placeholder="unbelivableemail@zmail.com" required />
	<input id="password" type="password" name="password" placeholder="123987" required />
	<input id="passwordDup" type="password" name="passwordDup" placeholder="123987" required />

	<p class="error">{errorMessage}</p>

	<button type="submit" name="submit"> SUBMIT </button>
</form>
