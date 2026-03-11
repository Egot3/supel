<script lang="ts">
	const { imageLink = '', caption, text } = $props();
</script>

/*got the joke? News are not countable*/
<div class="new">
	{#if imageLink}
		<img src={imageLink} alt="" class="img" />
	{/if}
	<div class="text-container">
		<p class="caption">
			{caption}
		</p>
		<div class="text">
			{text}
		</div>
	</div>
</div>

<style lang="sass">
    .caption 
        min-height: 1.5rem
        max-height: 5rem

    .text
        min-height: 3rem
        max-height: 30rem

    .img
        min-height: 4rem
        max-height: 10rem

    .new:nth-child(3n+1)
        grid-column: 1 / 4
    .new:nth-child(3n+2)
        grid-column: 5 / 8
    .new:nth-child(3n+3)
        grid-column: 9 / 12
</style>
