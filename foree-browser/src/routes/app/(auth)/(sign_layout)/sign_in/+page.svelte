<script lang="ts">
    import foreeLogo from "$lib/assets/images/foree_remittance_logo.svg"
</script>
<div>
    <header class="">
        <div class="logo"></div>
        <a href="/app/sign_up">Sign Up</a>
    </header>

    <main>
        <div class="left">
            <h2>
                Sigh-up & receive $20 to try the faster transfers to PPPPPP
            </h2>
            <ul>
                <li>&#10003; $0 fees and best FX rates</li>
                <li>&#10003; Transfers to 35+ banks</li>
                <li>&#10003; Cash pick-ups from 1500+ branches</li>
            </ul>
        </div>
        <div>

        </div>
    </main>
</div>

<style>
    header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 1rem 1rem;
    }

    header > a {
        text-decoration: none !important;
        font-size: 1em;
        font-weight: 600;
        background-color: var(--primary-color);
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

    main {

    }

    main > .left {
        color: var(--primary-color);
    }

    main > .left > h2 {
        margin-bottom: 1rem;
    }

    main > .left li {
        font-size: 1em;
        font-weight: 600;
    }

    main > .left li + li {
        margin-top: 0.5rem;
    }

    @media (max-width: 700px) {
        header {
            padding: 0.2rem 0.5rem;
        }
        .logo {
            width: 50px;
            height: 50px;
            background-image: url("$lib/assets/images/foree_remittance_small_logo.svg");
        }
    }
</style>