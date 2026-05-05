<script lang="ts">
	import '../app.css';
	let { children } = $props();
</script>

<div class="grid grid-cols-12 grid-rows-12 gap-gutter">
	{@render children()}
</div>
