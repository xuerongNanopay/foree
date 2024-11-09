<script lang="ts">
    const { children } = $props()
</script>


<nav class="header">
    Header
</nav>

<div class="main">
    <nav>
        <a href="dashboard">Dashboard</a>
    </nav>
    {@render children()}
</div>

<style>
    .header {
        z-index: 101;
        height: var(--foree-nav-height);
        padding: 0 var(--foree-page-padding-side);
        width: 100vw;
        margin: 0 auto;
        position:fixed;
        top:0;
        border: 1px solid red;

        @media (max-width: 831px) {
            & {
                top: unset;
                bottom: 0;
            }
        }
    }

    .main {
        padding-top: var(--foree-banner-height);
        height: 100%;
        border: 1px solid salmon;
        /* margin: o auto;
        position: relative; */

        @media (min-width: 832px) {
            & {
                padding-top: var(--foree-nav-height);
                padding-bottom: var(--foree-banner-height);
            }
        }
    }
</style>