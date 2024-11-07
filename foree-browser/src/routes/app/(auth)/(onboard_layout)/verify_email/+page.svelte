<script lang="ts">
    import { enhance } from "$app/forms"
    let submitting = $state(false)
    let inputPattern = "[0-9]{6}"

    let forgetPasswordForm = $state<VerifyEmailData>({
        code: "",
    })
    let forgetPasswordErr = $state<VerifyEmailError>({})

</script>

<main>
    <h2>Let's Verify Your Email Address</h2>
    <p>We have sent a verification code to your email. Please enter the code below to confirm that is account belongs to you.</p>

    <form
        method="POST"
        action="?/verify_code"
        use:enhance={
            () => {
                submitting = true
                return async ({update, result}) => {
                    await update()
                    submitting = false
                    if (result.type === "failure") {
                        //TODO
                        // forgetPasswordErr = {
                        //     ...result.data
                        // }
                    }
                }
            }
        }
    >
        <div>
            <label for="code">Verification Code</label>
            <input 
                type="text" 
                id="code" 
                name="code" 
                minlength=6
                maxlength=6
                pattern={inputPattern}
                placeholder="Please enter the 6-digit code"
                required
            />
        </div>
        <button disabled={submitting}>{submitting? "Submit..." : "Verify"}</button>
    </form>
    <form
        class="resend-code"
        method="POST"
        action="?/resend_code"
        use:enhance
    >
        <button>Resend Code</button>
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
    }

    form:first-of-type {
        margin: 2rem auto 0;
    }

    .resend-code {
        max-width: 500px;
        width: 100%;
        margin: 1rem auto;

        & button {
            display: block;
            outline: none;
            padding: 0.75rem 1rem;
            font-size: 1em;
            background: transparent;
            border-radius: 0.25rem;
            border: 1px solid var(--slate-600);
            color: var(--slate-600);
            margin: 0 auto;
        }
    }
</style>