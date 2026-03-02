<script lang="ts">
	import { applyAction, enhance } from '$app/forms';
	import Toaster from '$lib/components/Toaster.svelte';
	import '$lib/assets/styles/login.sass';
	import type { PageProps } from './$types';
	import { getContext } from 'svelte';

	const sendBread: (e: unknown) => void = getContext('sendBread');

	/** @type {import('./$types').PageProps} */
	let { form }: PageProps = $props();
</script>

<div class="login-header">LOGIN</div>
<form
	class="login-body"
	action="?/login"
	method="POST"
	use:enhance={({ formElement, formData, action, cancel }) => {
		return async ({ result }) => {
			await applyAction(result);

			if (form?.digitRequired) {
				sendBread('Digit is required');
			}
		};
	}}
>
	<input id="email" type="email" name="email" placeholder="unbelivableemail@zmail.com" required />
	<input id="password" type="password" name="password" placeholder="123987" required />

	<button type="submit"> SUBMIT </button>
</form>
