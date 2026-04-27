<script lang="ts">
	import { Modal, Button, Label, Input, Helper, Textarea, Fileupload } from 'flowbite-svelte';
	import { slide } from 'svelte/transition';
	import { UserEditOutline } from 'flowbite-svelte-icons';
	import type { User } from '$lib/types/user';
	import ImageFixer from '$lib/components/ImageFixer.svelte';

	let { open = $bindable(false), userInfo }: { open: boolean; userInfo: User } = $props();

	let files = $state<FileList | null>(null);

	let height = $state(0);
	let width = $state(0);
</script>

<Modal
	bind:open
	size="xs"
	class="bg-forest-900 dark:bg-forest-900 backdrop:bg-linen-900/50 text-linen-200 p-gutter mt-12 justify-self-center overflow-scroll leading-[-1rem]"
	transition={slide}
	dismissable={false}
>
	<form action="?/updateProfile" method="post" enctype="multipart/form-data">
		<h2 class="mb-5 text-lg font-normal flex text-light bg-dark">
			<UserEditOutline class="shrink-0 h-6 w-6 text-light bg-dark" />
			User updater
		</h2>

		<div class="text-left space-y-0.5 grid grid-cols-[65%_35%]">
			<div>
				<Label for="nickname">New nickname</Label>
				<Input
					defaultValue={userInfo.nickname}
					name="nickname"
					placeholder="extremly_cool_and_maybe_not_original_nickname103"
					class="rounded-lg bg-dark-compliment font-hollow text-[0.8rem] text-light pl-2"
				/>
				<Helper></Helper>

				<Label for="description">Description</Label>
				<Textarea
					required
					defaultValue={userInfo.description}
					name="description"
					placeholder="even cooler, originaler and longer text"
					class="rounded-lg w-full bg-dark-compliment font-hollow text-[0.8rem] text-light pl-2 resize-none"
					rows={5}
				/>
				<Helper class="bg-forest-900 text-linen-200 text-[0.75rem]">you can use MD there</Helper>

				<!-- {/* <Label for="image">Uploaded image</Label>
				<Fileupload
					type="image"
					name="image"
				>
				</Fileupload>
				<Helper>SVG, PNG, JPG or GIF (MAX. 800x400px).</Helper> * /} -->

				<Label for="with_helper" class="pb-2">Upload profile picture</Label>
				<Fileupload
					bind:files
					name="pfp"
					id="with_helper"
					class="rounded-lg bg-dark-compliment font-hollow text-[0.8rem] text-light pl-2"
					multiple={true}
					accept="image/*"
					clearable
				/>
				<Helper class="bg-forest-900 text-linen-200 text-[0.75rem]">WEBP, PNG, JPG.</Helper>
			</div>
			<div class="ml-10 mr-10">
				<ImageFixer bind:height bind:width bind:imageFiles={files} />
			</div>
		</div>

		<div class="space-x-2 bg-dark text-light mt-6">
			<Button
				type="submit"
				value="post"
				name="goPostIt"
				class="text-center font-medium inline-flex items-center justify-center text-linen-200 bg-coral-700 hover:bg-coral-800 dark:bg-coral-600 dark:hover:bg-coral-700 focus-within:ring-coral-300 dark:focus-within:ring-coral-900 px-5 py-2.5 text-sm focus-within:ring-4 focus-within:outline-hidden rounded-lg"
				color="red">Update NOW</Button
			>
			<Button
				//type="cancel"
				value="cancel"
				onclick={() => (open = false)}
				class="text-center font-medium inline-flex items-center justify-center text-linen-900 bg-transparent border border-linen-200 dark:border-linen-600 hover:bg-linen-100 dark:bg-linen-800 dark:text-linen-400 hover:text-primary-700 focus-within:text-primary-700 dark:focus-within:text-linen-50 dark:hover:text-linen-50 dark:hover:bg-linen-700 focus-within:ring-linen-200 dark:focus-within:ring-linen-700 px-5 py-2.5 text-sm focus-within:ring-4 focus-within:outline-hidden rounded-lg"
				color="dark">cancel</Button
			>
		</div>
	</form>
</Modal>
