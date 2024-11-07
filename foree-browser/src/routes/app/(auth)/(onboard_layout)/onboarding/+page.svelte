<script lang="ts">
    import { enhance } from "$app/forms"
    import MultiStepForm from "$lib/components/MultiStepForm.svelte"

    let countries = $state([
		{
			isoCode:  "CA",
			name:  "Canada",
		},
		{
			isoCode:  "PK",
			name:  "Pakistan",
		},
		{
			isoCode:  "US",
			name:  "United States of America",
		},
	]);

    let submitting = $state(false)
    let createUserForm = $state<CreateUserData>({
        firstName: "",
        middleName: "",
        lastName: "",
        dob: "",
        pob: "",
        nationality: "",
        address1: "",
        address2: "",
        city: "",
        province: "",
        country: "",
        postalCode: "",
        phoneNumber: "",
        identificationType: "",
        identificationValue: "",
    })
    let createUserErr = $state<CreateUserError>({})

    $inspect(createUserForm)
</script>

{#snippet step1()}
    <div>
    </div>
{/snippet}

{#snippet step2()}
    <div class="step">
        <h2>Your Residential Address and Phone Number</h2>
        <p>We require this information to continue setting up your Foree Remittance account.</p>
        <div class="fields">
            <div class="country">
                <label for="country">Country</label>
                <!-- <input bind:value={createUserForm.country} type="country" id="country" name="country" required> -->
                 <select
                    bind:value={createUserForm.country}
                    id="country"
                    name="country"
                    required
                 >
                    {#each countries as country}
                        <option value={country.isoCode}>{country.name}</option>
                    {/each}
                 </select>
                {#if !!createUserErr?.country}
                    <p class="input-error">{createUserErr.country}</p>
                {/if}
            </div>
            <div class="address1">
                <label for="address1">Address</label>
                <input bind:value={createUserForm.address1} type="text" id="address1" name="address1" required>
                {#if !!createUserErr?.address1}
                    <p class="input-error">{createUserErr.address1}</p>
                {/if}
            </div>
            <div class="address2">
                <label for="address2">Address Line2(Apt,suite,etc.)</label>
                <input bind:value={createUserForm.address2} type="text" id="address2" name="address2">
                {#if !!createUserErr?.address2}
                    <p class="input-error">{createUserErr.address2}</p>
                {/if}
            </div>
        </div>
    </div>
{/snippet}

{#snippet step3()}
    <div>
    </div>
{/snippet}
    
<main>
    <MultiStepForm steps={[step1, step2, step3]}>
    </MultiStepForm>
</main>

<style>
    main {
        width: 95%;
        max-width: 900px;
        margin: 3.5rem auto;
        border: 1px solid salmon;
    }

    .step {
        & > p:first-of-type {
            text-align: center;
            margin: 1.5rem 0;
        }

        & > h2:first-of-type {
            text-align: center;
            margin: 0;
        }

        @media (max-width: 400px) {
            & > h2:first-of-type {
                font-size: 1em;
            }
        }
    }

    .step {
        overflow: hidden;
        text-overflow: ellipsis; 

        & label {
            display: block;
            color: var(--slate-500);
            margin-bottom: 0.25rem;
        }

        & input, & select {
            display: block;
            width: 100%;
            border: 1px solid var(--slate-400);
            color: var(--slate-600);
            font-size: 1.25em;
            padding: 0.5rem 0.25rem;
            border-radius: 0.25rem;

            &:focus {
                outline: none !important;
                border-color: var(--primary-color);
                box-shadow: 0 0 5px var(--primary-color);
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

        .input-error {
            color: var(--rose-400);
        }
    }

    .fields {
        display: grid;
        grid-template-columns: repeat(12, 1fr);
        gap: 1rem;

        .country {
            grid-column: 1 / span 12;
        }

        .address1 {
            grid-column: 1 / span 8;
        }

        .address2 {
            grid-column: 9 / 13;
        }

        @media (max-width: 700px) {
            .address1 {
                grid-column: 1 / span 12;
            }

            .address2 {
                grid-column: 1 / span 12;
            }
        }
    }

</style>