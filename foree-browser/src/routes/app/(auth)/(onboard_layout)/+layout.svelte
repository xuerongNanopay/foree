<script lang="ts">
    const { children } = $props()
</script>

<div class="wrapper">
    <header>
        <div class="logo">

        </div>
    </header>
    <main>
        {@render children()}
    </main>
</div>

<style>
    .wrapper {
        height: auto;
        min-height: 100vh;
    }
    header {
        position: relative;
        border-bottom: 2px solid var(--emerald-800);
        height: 70px;
    }
    header .logo {
        /* height: 100%; */
        position: absolute;
        height: 60px;
        width: 120px;
        left: 50%;
        top: 50%;
        translate: -50% -50%;
        background-size: 100% 100%;
        background-image: url("$lib/assets/images/foree_remittance_logo.svg");
    }
</style>