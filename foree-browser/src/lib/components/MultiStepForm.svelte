<script lang="ts">
    import type { Snippet } from "svelte"
    interface Props {
        steps: Snippet<[]>[];
    }

    let { steps }: Props = $props()
</script>

<form>
    <div style:width={steps.length*100 + "%"}>
        {#each steps as step}
            <div>
                {@render step()}
            </div>
        {/each}
    </div>
</form>


<style>
    form > div {
        display: flex;
        & > div {
        width: 100%;
       } 
    }
</style>