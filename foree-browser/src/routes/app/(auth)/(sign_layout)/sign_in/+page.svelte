<script lang="ts">
    import foreeLogo from "$lib/assets/images/foree_remittance_logo.svg"
</script>
<div>
    <header class="">
        <div class="logo"></div>
        <a href="/app/sign_up">Sign Up</a>
    </header>

</div>

<style>
    header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 1rem 1rem;
    }

    a {
        text-decoration: none !important;
        font-size: 1em;
        font-weight: 600;
        background-color: #1A6B54;
        padding: 0.6rem 0.8rem;
        color: #fff;
        border-radius: 5px;
    }

    .logo {
        width: 100px;
        height: 40px;
        background-size: 100% 100%;
        background-image: url("$lib/assets/images/foree_remittance_logo.svg");
    }

    @media (max-width: 700px) {
        .logo {
            width: 50px;
            height: 50px;
            background-image: url("$lib/assets/images/foree_remittance_small_logo.svg");
        }
    }
</style>