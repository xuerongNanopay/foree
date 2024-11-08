<script lang="ts">
    import { enhance } from "$app/forms"
    let submitting = $state(false)
    let forgetPasswordForm = $state<ForgetPasswordData>({
        email: "",
    })
    let forgetPasswordErr = $state<ForgetPasswordError>({})
    let dialog: HTMLDialogElement
</script>

<dialog bind:this={dialog}>
    <h2>Notice</h2>
    <p>We will send you email. Please follow the instruction in the email.</p>
    <button type="button" onclick={() => dialog.close()}>OK</button>
</dialog>

<main>
    <h2>Forget Password?</h2>
    <p>Enter email you used to create your account in order to reset your password.</p>
    
    <form 
        method="POST" 
        use:enhance={
            () => {
                submitting = true
                return async ({update, result}) => {
                    await update()
                    submitting = false
                    forgetPasswordErr={}
                    if (result.type === "success") {
                        dialog.showModal()
                    } else if (result.type === "failure") {
                        forgetPasswordErr = {
                            ...result.data
                        }
                    }
                }
            }
        }
    >
        <div>
            <label for="email">Email</label>
            <input bind:value={forgetPasswordForm.email} type="email" id="email" name="email" required>
            {#if !!forgetPasswordErr?.email}
                <p class="input-error">{forgetPasswordErr.email}</p>
            {/if}
        </div>
        <button disabled={submitting}>{submitting? "Submit..." : "Submit"}</button>
        {#if !!forgetPasswordErr?.cause}
            <p class="input-error">{forgetPasswordErr.cause}</p>
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

        & > button {
            margin-top: 1.5rem;
        }
    }

    dialog {
        top: 50%;
        left: 50%;
        translate: -50% -50%;
        width: 80%;
        max-width: 400px;
        outline: none;
        border: 1px solid var(--primary-color);
        border-radius: 7px;
        padding: 0.75rem 0.5rem;

        & h2 {
            color: var(--primary-color);
            border-bottom: 1px solid var(--primary-color);
        }

        & p {
            color: var(--primary-color);
            margin: 1rem 0.25rem;
        }

        & button {
            display: block;
            background: transparent;
            border: 1px solid var(--primary-color);
            padding: 0.5rem 1rem;
            color: var(--primary-color);
            border-radius: 8px;
            font-size: large;
            font-weight: 600;
            width: 6rem;
            margin: 0 auto;
            transition: background-color 0.5s ease-in-out;

            &:hover {
                background-color: var(--emerald-100);
            }
        }

        &::backdrop {
            background-color: var(--slate-400);
            opacity: 0.5;
        }
    }


</style>