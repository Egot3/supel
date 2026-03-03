<script lang="ts">
	import Toast from './Toast.svelte';
	let { messages = $bindable(new Map<string, unknown>()) } = $props();
	function RemoveMessage(key: string) {
		messages = new Map([...messages].filter(([k]) => k !== key)); //don't ask
	}
</script>

<div>
	{#each messages as [key, message] (key)}
		<Toast cooked={() => RemoveMessage(key)}>
			{message}
		</Toast>
	{/each}
</div>
