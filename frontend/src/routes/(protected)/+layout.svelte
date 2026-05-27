<script lang="ts">
	import { getUserContext, setUserContext } from '$lib/context/userContext';
	import type { User } from '$lib/types/user';
	import type { LayoutData } from './$types';
	import type { Snippet } from 'svelte';

	let { data, children }: { data: LayoutData; children: Snippet } = $props();

	let user = $state(((data.user as User) ?? {}) as User);
	try {
		/* console.log('user and ctx: ', user, getUserContext()); */

		if (user != getUserContext()) {
			setUserContext(user);
		}
	} catch {
		setUserContext(user);
	}

	/* $inspect(user); */
</script>

<a
	href={`http://localhost:5173/profile/${user.uuid}`}
	aria-label="Your profile"
	class="col-start-12 col-end-13 bg-forest-900 rounded-xl p-1 aspect-square"
>
	<img src={user.avatarUrl} alt="your profile pic" class="rounded-full aspect-square" />
</a>
{@render children()}
