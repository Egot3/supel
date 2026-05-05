<script lang="ts">
	import { setUserContext } from '$lib/context/userContext';
	import type { User } from '$lib/types/user';
	import type { LayoutData } from './$types';
	import type { Snippet } from 'svelte';
	import '../app.css';

	let { data, children }: { data: LayoutData; children: Snippet } = $props();

	let user = $state((data.user as User) ?? ({} as User));
	setUserContext(user);
	$inspect(user);
</script>

<div class="grid grid-cols-12 grid-rows-12 gap-gutter">
	<a
		href="http://localhost:5173/profile"
		aria-label="Your profile"
		class="col-start-12 col-end-13 bg-forest-900 rounded-xl p-1"
	>
		<img src={user.avatarUrl} alt="your profile pic" class="rounded-full" />
	</a>
	{@render children()}
</div>
