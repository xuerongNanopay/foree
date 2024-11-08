<script lang="ts">
    import { enhance } from "$app/forms"
    import MultiStepForm from "$lib/components/MultiStepForm.svelte"

    const phoneNumberPattern = "[0-9]{3}-[0-9]{2}-[0-9]{3}"

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
    <div class="step">
        <h2>Let's Get to Know You!</h2>
        <p>Please enter your full legal name se we can begin setting up your account.</p>
        <div class="name-fields">
            <div class="firstname">
                <label for="firstName">First Name</label>
                <input bind:value={createUserForm.firstName} type="text" id="firstName" name="firstName" required>
                {#if !!createUserErr?.firstName}
                    <p class="input-error">{createUserErr.firstName}</p>
                {/if}
            </div>
            <div class="middleName">
                <label for="middleName">Middle Name</label>
                <input bind:value={createUserForm.middleName} type="text" id="middleName" name="middleName">
                {#if !!createUserErr?.middleName}
                    <p class="input-error">{createUserErr.middleName}</p>
                {/if}
            </div>
            <div class="lastName">
                <label for="lastName">Last Name</label>
                <input bind:value={createUserForm.lastName} type="text" id="lastName" name="lastName" required>
                {#if !!createUserErr?.lastName}
                    <p class="input-error">{createUserErr.lastName}</p>
                {/if}
            </div>
        </div>
    </div>
{/snippet}

{#snippet step2()}
    <div class="step">
        <h2>Your Residential Address and Phone Number</h2>
        <p>We require this information to continue setting up your Foree Remittance account.</p>
        <div class="address-fields">
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
            <div class="city">
                <label for="city">City</label>
                <input bind:value={createUserForm.city} type="text" id="city" name="city">
                {#if !!createUserErr?.city}
                    <p class="input-error">{createUserErr.city}</p>
                {/if}
            </div>
            <div class="province">
                <label for="province">province</label>
                 <select
                    bind:value={createUserForm.province}
                    id="province"
                    name="province"
                    required
                 >
                    {#each countries as province}
                        <option value={province.isoCode}>{province.name}</option>
                    {/each}
                 </select>
                {#if !!createUserErr?.province}
                    <p class="input-error">{createUserErr.province}</p>
                {/if}
            </div>
            <div class="postal-code">
                <label for="postalCode">Postal Code</label>
                <input bind:value={createUserForm.postalCode} type="text" id="postalCode" name="postalCode">
                {#if !!createUserErr?.postalCode}
                    <p class="input-error">{createUserErr.postalCode}</p>
                {/if}
            </div>
            <div class="phone-number">
                <label for="phoneNumber">Phone Number</label>
                <input 
                    bind:value={createUserForm.phoneNumber} 
                    type="text" 
                    id="phoneNumber" 
                    name="phoneNumber"
                    placeholder="000-000-0000"
                    pattern={phoneNumberPattern}
                >
                {#if !!createUserErr?.phoneNumber}
                    <p class="input-error">{createUserErr.phoneNumber}</p>
                {/if}
            </div>
        </div>
    </div>
{/snippet}

{#snippet step3()}
    <div class="step">
        <h2>Personal Details</h2>
        <p>Almost done! Infomation below is requested by XXXXXX Bank of XXXXXXX, our Foree Remittance payout partner, in order to process your transfers under Pkakistani regulatory guidelines</p>
        <div class="personal-fields">
            <div class="dob">
                <label for="dob">Date of Birth</label>
                <input 
                    bind:value={createUserForm.dob}
                    type="date" 
                    id="dob" 
                    name="dob"
                    required
                >
                {#if !!createUserErr?.dob}
                    <p class="input-error">{createUserErr.dob}</p>
                {/if}
            </div>
            <div class="pob">
                <label for="country">Place of Birth</label>
                <!-- <input bind:value={createUserForm.country} type="country" id="country" name="country" required> -->
                 <select
                    bind:value={createUserForm.pob}
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
            <div class="nationality">
                <label for="nationality">Nationality</label>
                <!-- <input bind:value={createUserForm.nationality} type="nationality" id="nationality" name="nationality" required> -->
                 <select
                    bind:value={createUserForm.nationality}
                    id="nationality"
                    name="nationality"
                    required
                 >
                    {#each countries as nationality}
                        <option value={nationality.isoCode}>{nationality.name}</option>
                    {/each}
                 </select>
                {#if !!createUserErr?.nationality}
                    <p class="input-error">{createUserErr.nationality}</p>
                {/if}
            </div>
            <div class="identification-type">
                <label for="identificationType">Identification Document Type</label>
                <!-- <input bind:value={createUserForm.identificationType} type="identificationType" id="identificationType" name="identificationType" required> -->
                 <select
                    bind:value={createUserForm.identificationType}
                    id="identificationType"
                    name="identificationType"
                    required
                 >
                    {#each countries as identificationType}
                        <option value={identificationType.isoCode}>{identificationType.name}</option>
                    {/each}
                 </select>
                {#if !!createUserErr?.identificationType}
                    <p class="input-error">{createUserErr.identificationType}</p>
                {/if}
            </div>
            <div class="identification-value">
                <label for="identificationValue">Identification Number</label>
                <input 
                    bind:value={createUserForm.identificationValue}
                    id="identificationValue" 
                    name="identificationValue"
                    required
                >
                {#if !!createUserErr?.identificationValue}
                    <p class="input-error">{createUserErr.identificationValue}</p>
                {/if}
            </div>
        </div>
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

        @media (max-width: 700px) {
            margin: 2.5rem auto;
        }
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

    .name-fields{
        display: grid;
        gap: 1rem;
        grid-template-columns: repeat(3, 1fr);
        @media (max-width: 700px) {
            & {
                grid-template-columns: 1fr;
            }
        }
    }

    .address-fields {
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

        .city {
            grid-column: 1 / 5;
        }

        .province {
            grid-column: 5 / 9;
        }

        .postal-code {
            grid-column: 9 / 13;
        }

        .phone-number {
            grid-column: 1 / span 12;
        }

        @media (max-width: 700px) {
            .address1 {
                grid-column: 1 / span 12;
            }

            .address2 {
                grid-column: 1 / span 12;
            }
            .city {
                grid-column: 1 / span 12;
            }

            .province {
                grid-column: 1 / span 12;
            }

            .postal-code {
                grid-column: 1 / span 12;
            }
        }
    }

    .personal-fields {
        display: grid;
        gap: 1rem;
    }
</style>