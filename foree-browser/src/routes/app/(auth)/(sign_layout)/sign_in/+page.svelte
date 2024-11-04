<script lang="ts">
    import foreeLogo from "$lib/assets/images/foree_remittance_logo.svg"
</script>
<div class="contain">
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
        <div class="right">
            aaa
        </div>
    </main>
</div>

<style>
    .contain {
        display: flex;
        flex-direction: column;
        min-height: 100vh;
    }
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
        flex-grow: 1;
        display: grid;
        grid-template-columns: 1fr 1fr;
        align-items: center;
        justify-items: center;
    }

    main > .left {
        color: var(--primary-color);
        padding: 0 1rem;
    }

    main > .left > h2 {
        margin-bottom: 1rem;
        font-size: clamp(1.4em, 2.5vw, 6em);
    }

    main > .left > ul {
        list-style-type: none;
    }

    main > .left li {
        font-size: clamp(0.9em, 2.5vw, 1.1em);
        font-weight: 600;
    }

    main > .left li + li {
        margin-top: 0.5rem;
    }

    @media (max-width: 956px) {
        header {
            padding: 0.2rem 0.5rem;
        }
        .logo {
            width: 50px;
            height: 50px;
            background-image: url("$lib/assets/images/foree_remittance_small_logo.svg");
        }
        main {
            grid-template-columns: 1fr;
            grid-template-rows: min-content 1fr;
            align-items: normal;
        }
        main > .right {
            margin-top: 1rem;
        }
    }
</style>