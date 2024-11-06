<script lang="ts">
    import { page } from '$app/stores'
    import { enhance } from "$app/forms"
    import eyeIcon from "$lib/assets/icons/eye.png"
    import eyeHideIcon from "$lib/assets/icons/eye_hide.png"

    let isHidePassword = $state(true)
    let submitting = $state(false)

    let renewPasswordForm = $state<RenewPasswordData>({
        password: "",
        rePassword: "",
        retrieveCode: $page.url.searchParams.get('retrieveCode') ?? ""
    })

    let renewPasswordErr = $state<RenewPasswordError>({})

    function toggleEye() {
        isHidePassword = !isHidePassword
    }
</script>

<main>
    <h2>Renew Password</h2>
    <p>Please enter a new password.</p>
    <form
        method="POST"
        use:enhance={
            () => {
                submitting = true
                renewPasswordErr={}
                return async ({update, result}) => {
                    await update()
                    submitting = false
                    if (result.type === "failure") {
                        renewPasswordErr = {
                            ...result.data
                        }
                    }
                }
            }
        }
    >
        <div class="password">
            <label for="password">Password</label>
            <div>
                <input bind:value={renewPasswordForm.password} type={isHidePassword ? "password" : "text"} id="password" name="password" required>
                <button type="button" onclick={toggleEye}>
                    <img src={isHidePassword ? eyeHideIcon : eyeIcon} alt=""/>
                </button>
            </div>
            {#if !!renewPasswordErr?.password}
                <p class="input-error">{renewPasswordErr.password}</p>
            {/if}
        </div>
        <div class="password">
            <label for="rePassword">Retype Password</label>
            <div>
                <input bind:value={renewPasswordForm.rePassword} type={isHidePassword ? "password" : "text"} id="rePassword" name="rePassword" required>
                <button type="button" onclick={toggleEye}>
                    <img src={isHidePassword ? eyeHideIcon : eyeIcon} alt=""/>
                </button>
            </div>
            {#if !!renewPasswordErr?.rePassword}
                <p class="input-error">{renewPasswordErr.rePassword}</p>
            {/if}
        </div>
        <button disabled={submitting}>{submitting? "Submit..." : "Update"}</button>
        {#if !!renewPasswordErr?.cause}
            <p class="input-error">{renewPasswordErr.cause}</p>
        {/if}
    </form>
</main>

<style>
    main {
        width: 95%;
        max-width: 700px;
        margin: 0 auto;
    }

    main > :is(h2, p) {
        margin-top: 2.5rem;
        text-align: center;
    }

    main > p {
        margin-top: 2rem;
        color: var(--slate-800)
    }

    main > form {
        max-width: 500px;
        width: 100%;
        margin: 2rem auto;
    }

    form {
        & input {
            display: block;
            width: 100%;
            border: 1px solid var(--slate-400);
            color: var(--slate-600);
            font-size: 1.25em;
            padding: 0.5rem 0.25rem;
            border-radius: 0.25rem;

            &:focus {
                outline: none !important;
                border-color: var(--emerald-800);
                box-shadow: 0 0 5px var(--emerald-800);
            }
        }
        & label {
            display: block;
            color: var(--slate-500);
            margin-bottom: 0.25rem;
        }

        & > div:nth-child(2){
            margin-top: 1.25rem;
        }

        .password div {
            position: relative;

            & button {
                position: absolute;
                right: 0.75rem;
                top: 50%;
                transform: translateY(-50%);
                height: 1.5rem;
                width: 1.5rem;
                border: 0;
                background: transparent;

                & img {
                    height: 100%;
                    width: 100%;
                }
            }
        }

        & > button {
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

            &:hover {
                transform: scale(1.01);
            }
        }
        .input-error {
            color: var(--rose-400);
        }
    }
</style>