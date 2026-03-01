<script lang="ts">
	import { applyAction, enhance } from '$app/forms';
    import '$lib/assets/styles/login.sass'
	import type { PageProps } from './$types';

    /** @type {import('./$types').PageProps} */
	let { form }:PageProps = $props();
</script>

<div style="display: grid; grid-template: repeat(12, 1fr) / repeat(12, 1fr); height: 100vh; width:100vw">
    <div class="login-header">
        LOGIN
    </div>
    <form class="login-body" action="?/login" method="POST" use:enhance={({ formElement, formData, action, cancel })=>{
        return async({ result })=>{
            await applyAction(result)
        }
    }}>
        {#if form?.missing}<p class="error">The email and password fields are required</p>{/if}
        {#if form?.digitRequired}<p class="error">Password must have at least 4 digits</p>{/if}
        {#if form?.whitespace}<p class="error">Password mustn't have any space character</p>{/if}
        {#if form?.upperCaseRequired}<p class="error">Password must have at least one UPPERCASE letter</p>{/if}
        {#if form?.lowerCaseRequired}<p class="error">Password must have at least one lowercase letter</p>{/if}

        <input id="email" type="email" name="email" placeholder="unbelivableemail@zmail.com" required/>
        <input id="password" type="password" name="password" placeholder="123987" required/>
        
        <button type="submit">
            SUBMIT
        </button>
    </form>
</div>
