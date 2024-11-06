<script lang="ts">

    let submitting = $state(false)
    let forgetPasswordForm = $state<ForgetPasswordData>({
        email: "",
    })
    let forgetPasswordErr = $state<ForgetPasswordError>({})
    let dialog: HTMLDialogElement
    $effect(() => {
		// dialog.showModal()
	})
</script>

<dialog bind:this={dialog}>
    <h2>Info</h2>
    <p>Message</p>
    <button type="button">OK</button>
</dialog>

<main>
    <h2>Forget Password?</h2>
    <p>Enter email you used to create your account in order to reset your password.</p>
    
    <form method="POST">
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
        max-width: 900px;
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
        margin-top: 2rem;
    }

    form {
        & label {
            display: block;
            color: var(--slate-500);
            margin-bottom: 0.25rem;
        }

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

        & button {
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

    }


    dialog {
        top: 50%;
        left: 50%;
        translate: -50% -50%;
        width: 80%;
        max-width: 400px;
        outline: none;
        border: 1px solid var(--emerald-800);
        border-radius: 7px;
        padding: 0.75rem 0.5rem;
    }

    dialog::backdrop {
        background-color: var(--slate-400);
        opacity: 0.5;
    }



</style>