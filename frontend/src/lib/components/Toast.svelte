<script lang="ts">
	import { onDestroy } from 'svelte';

	let { cooked, children } = $props();

	const timer = setTimeout(() => {
		cooked();
	}, 2000);

	onDestroy(() => clearTimeout(timer));
</script>

<button class="bread" onclick={() => cooked()}>
	{@render children()}
</button>

<style lang="sass">

</style>
