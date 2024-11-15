<script lang="ts">
    import xIcon  from "$lib/assets/icons/x.png"
    import { afterNavigate } from "$app/navigation";
    const { children } = $props()

    afterNavigate(({from})=> {
        console.log(from)
    })
</script>

<div class="wrapper">
    <header>
        <div class="logo"></div>
        <button>
            <img src={xIcon} alt="close"/>
        </button>
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

    .wrapper > header {
        position: sticky;
        top: 0px;
        border-bottom: 2px solid var(--emerald-800);
        height: 70px;
        background: white;
    }

    .wrapper > header .logo {
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

    .wrapper > header > button {
        outline: none;
        border: none;
        background: transparent;
    }

    .wrapper > header > button img {
        height: 18px;
        width: 18px;
        position: absolute;
        top: 50%;
        right: 1.25rem;
        translate: 0 -50%;
        transition: transform .25s ease-in-out;
        &:hover {
            transform: scale(1.25);
        }
    }

    @media (max-width: 956px) {
        .wrapper > header {
            height: 60px;
        }
        .wrapper > header .logo {
            width: 50px;
            height: 50px;
            background-image: url("$lib/assets/images/foree_remittance_small_logo.svg");
        }
        .wrapper > header form img {
            right: 1rem;
        }
    }
</style>