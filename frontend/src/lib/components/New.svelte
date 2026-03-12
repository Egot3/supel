<script lang="ts">
	const { imageLink = '', caption, text } = $props();
</script>

<!-- /*got the joke? News are not countable*/ -->
<div class="new">
	{#if imageLink}
		<img src={imageLink} alt="" class="img" />
	{/if}
	<div class="text-container">
		<p class="caption">
			{caption}
		</p>
		<br />
		<div class="text">
			{text}
		</div>
	</div>
</div>

<style lang="sass">
    @use '../assets/_variables' as *
    @use 'sass:math'
    .new
        background: $blacker
        border: $whiter 0.5rem solid
        padding: 1rem
        border-radius: 1rem
        
    .caption 
        min-height: 1.5rem
        max-height: 5rem

    .text
        min-height: 2rem
        max-height: 30rem
        overflow: scroll
        -ms-overflow-style: none
        scrollbar-width: none  //--ux
        margin-bottom: 0.3rem
        $start-opaque: 80%
        mask-image: linear-gradient(to bottom, 
            black 0%, 
            black $start-opaque, 
            transparent math.min(100%, $start-opaque + 20%)
        )
    
    .text::-webkit-scrollbar //PWA
        display: none
        width: 0
        height: 0

    .img
        min-height: 4rem
        max-height: 10rem

    // .new:nth-child(3n+1)
    //     grid-column: 1 / 5
    // .new:nth-child(3n+2)
    //     grid-column: 5 / 9
    // .new:nth-child(3n+3)
    //     grid-column: 9 / 13
</style>
