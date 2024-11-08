<script lang="ts">
    import type { Snippet } from "svelte"
    interface Props {
        steps: Snippet<[]>[];
    }

    let { steps }: Props = $props()
</script>

<form style:width={steps.length*100 + "%"}>
    {#each steps as step}
        <div>
            {@render step()}
        </div>
    {/each}
</form>


<style>
    form {
        display: flex;
        & > div {
        width: 100%;
       } 
    }
</style>