<script lang="ts">
    import type { Snippet } from "svelte"
    interface Props {
        steps: Snippet<[]>[];
    }

    let { steps }: Props = $props()
</script>

<form>
    {#each steps as step}
        {@render step()}
    {/each}
</form>


<style>

</style>