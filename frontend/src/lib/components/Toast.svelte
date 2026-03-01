<script lang="ts">
    import { onDestroy } from "svelte";
	
    let { cooked, children } = $props()

    const timer = setTimeout(()=>{
        cooked(false)
    })
</script>

<button class="Bread" onclick={()=> cooked(true)}>
    {@render children()}
</button>

<style lang="sass">

</style>