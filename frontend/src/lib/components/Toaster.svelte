<script lang="ts">
	import Toast from './Toast.svelte';
	let { messages = $bindable([]) } = $props();
</script>

<div>
	{#each messages as message, index}
		<Toast
			cooked={(cooked: boolean) => {
				if (cooked) {
					messages.splice(index, 1);
				}
			}}
		>
			{message}
		</Toast>
	{/each}
</div>
