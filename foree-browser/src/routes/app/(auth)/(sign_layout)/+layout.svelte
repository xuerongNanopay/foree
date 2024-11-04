<script lang="ts">
    const { children } = $props()
</script>

<div>
    {@render children()}
</div>

<style>
    div {
        height: auto;
        min-height: 100vh;
        background: url("$lib/assets/images/foree_background.svg")
    }
</style>