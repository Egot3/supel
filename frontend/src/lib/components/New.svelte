<script lang="ts">
	import { Card } from 'flowbite-svelte';
	const { imageLink = '', caption, text } = $props();
</script>

<!-- /*got the joke? News are not countable*/ -->
<Card img={imageLink} class="mb-2 break-inside-avoid">
	<div class="m-6">
		<h5 class="mb-2 text-2xl font-bold bg-dark text-light font-hollow min-h-6 max-h-20">
			{caption}
		</h5>
		<p
			class="mb-3 leading-tight font-normal bg-dark text-light font-hollow min-h-8 max-h-80 overflow-scroll no-scrollbar mask-[linear-gradient(to_bottom,black_0%,black_80%,transparent_100%)]"
		>
			{text}
		</p>
	</div>
</Card>

<!-- <style lang="sass">
    @use '../assets/_variables' as *
    @use 'sass:math'
    .new
        //margin: 0.5rem
        background: $blacker
        border: $whiter 0.5rem solid
        padding: 1rem
        border-radius: 1rem
        
    .caption 
        min-height: 1.5rem
        max-height: 5rem

    .text
        min-height: 2rem
        max-height: 10rem
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
</style> -->
