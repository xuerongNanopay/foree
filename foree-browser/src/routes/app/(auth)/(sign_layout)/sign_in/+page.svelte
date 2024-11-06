<script lang="ts">
    import { enhance } from "$app/forms"
    import appStoreBadge from "$lib/assets/images/app_store_badge.svg"
    import playStoreBadge from "$lib/assets/images/play_store_badge.svg"
    import eyeIcon from "$lib/assets/icons/eye.png"
    import eyeHideIcon from "$lib/assets/icons/eye_hide.png"

    let isHidePassword = $state(true)
    let submitting = $state(false)

    let signInForm = $state<SignInFormData>({
        email: "",
        password: ""
    })

    let signInErr = $state<SignInFormError>({})

    function toggleEye() {
        isHidePassword = !isHidePassword
    }
    $inspect(signInErr)

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
        <div class="sign-in">
            <h3>Welcome Back</h3>
            <form 
                method="POST" 
                use:enhance={
                    () => {
                        submitting = true
                        signInErr = {}
                        return async ({update, result}) => {
                            await update()
                            submitting = false
                            if (result.type === "failure") {
                                signInErr = {
                                    ...result.data
                                }
                            }
                        }
                    }
                }
            >
                <div class="email">
                    <label for="email">Email</label>
                    <input bind:value={signInForm.email} type="email" id="email" name="email" required>
                    {#if !!signInErr?.email}
                        <p class="input-error">{signInErr.email}</p>
                    {/if}
                </div>
                <div class="password">
                    <label for="password">Password</label>
                    <div>
                        <input bind:value={signInForm.password} type={isHidePassword ? "password" : "text"} id="password" name="password" required>
                        <button type="button" onclick={toggleEye}>
                            <img src={isHidePassword ? eyeHideIcon : eyeIcon} alt=""/>
                        </button>
                    </div>
                    {#if !!signInErr?.password}
                        <p class="input-error">{signInErr.password}</p>
                    {/if}
                </div>
                <button disabled={submitting}>{submitting? "Submit..." : "Sign In"}</button>
                {#if !!signInErr?.cause}
                    <p class="input-error">{signInErr.cause}</p>
                {/if}
            </form>
            <a class="forget-password" href="forget_password">Forget Password?</a>
            <div class="mobile-badge">
                <a href="http://www.google.ca">
                    <img src={appStoreBadge} alt="App Store"/>
                </a>
                <a href="http://www.google.ca">
                    <img src={playStoreBadge} alt="Play Store"/>
                </a>                
            </div>
            <div class="mobile-badge-copyright">
                <p>Apple and the Apple Logo are trademarks of Apple Inc.</p>
                <p>Google Play and the Google Play logo are trademarks of Google LLC.</p>
            </div>
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
        transition: transform .25s ease-in-out;
    }

    header > a:hover {
        transform: scale(1.05);
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
        font-size: clamp(1em, 2.5vw, 1.2em);
        font-weight: 600;
    }

    main > .left li + li {
        margin-top: 0.5rem;
    }

    .sign-in {
        width: 90%;
        background-color: #fff;
        border-radius: 2rem;
        padding: 2rem 1rem;
        box-shadow: 0px 5px 5px 2px rgba(0, 0, 0, 0.2);
    }

    .sign-in > form {
        margin-top: 1rem;
    }
    
    /* .sign-in form div:not(:first-child) {
       margin-top: 1.25rem;
    } */

    .sign-in form > div:nth-child(2){
       margin-top: 1.25rem;
    }

    .sign-in label {
        display: block;
        color: var(--slate-500);
        margin-bottom: 0.25rem;
    }

    .sign-in input {
        display: block;
        width: 100%;
        border: 1px solid var(--slate-400);
        color: var(--slate-600);
        font-size: 1.25em;
        padding: 0.5rem 0.25rem;
        border-radius: 0.25rem;
    }

    .sign-in input:focus { 
        outline: none !important;
        border-color: var(--emerald-800);
        box-shadow: 0 0 5px var(--emerald-800);
    }

    .sign-in > form .password div {
        position: relative;
    }

    .sign-in > form .password div button {
        position: absolute;
        right: 0.75rem;
        top: 50%;
        transform: translateY(-50%);
        height: 1.5rem;
        width: 1.5rem;
        border: 0;
        background-color: white;
    }

    .sign-in > form .password div button img {
        height: 100%;
        width: 100%;
    }

    .sign-in .forget-password {
        display: inline-block;
        text-decoration: none !important;
        color: var( --slate-400);
        margin-top: 0.5rem;
    }

    .sign-in > form > button {
        display: block;
        width: 100%;
        margin-top: 1.5rem;
        background-color: var(--primary-color);
        border: 0px;
        padding: 0.75rem 0;
        border-radius: 0.25rem;
        color: white;
        font-size: 1em;
        font-weight: 600;
        transition: transform .25s ease-in-out;
    }

    .sign-in > form > button:hover {
        transform: scale(1.01);
    }

    .sign-in > form > button:disabled,button[disabled] {
        opacity: 0.6;
    }

    .mobile-badge {
        display: flex;
        gap: 1rem;
        justify-content: center;
        height: 2rem;
        margin-top: 1.25rem;
        width: 100%;
    }

    .mobile-badge a {
        display: inline-block;
        height: 100%;
    }
    .mobile-badge img {
        display: inline-block;
        height: 100%;
    }

    .mobile-badge-copyright {
        margin-top: 0.5rem;
        text-align: center;
        color: var(--slate-600);
        font-size: 0.5rem;
    }

    .sign-in form .input-error {
        color: var(--rose-400);
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
        .sign-in {
            margin-top: 1rem;
            width: 100%;
            border-radius: 2rem 2rem 0 0;
            box-shadow: 0;
        }
    }
</style>