<script lang="ts">
	import { setUserContext } from '$lib/context/userContext';
	import type { User } from '$lib/types/user';
	import '../app.css';
	import type { LayoutData } from './$types';
	import type { Snippet } from 'svelte';

	let { data, children }: { data: LayoutData; children: Snippet } = $props();

	let user = $state((data.user as User) ?? ({} as User));
	setUserContext(user);
	$inspect(user);
</script>

<div class="grid grid-cols-12 grid-rows-12 h-screen w-screen">
	{@render children()}
</div>
