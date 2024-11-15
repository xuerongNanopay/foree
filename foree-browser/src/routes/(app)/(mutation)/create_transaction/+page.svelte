<script lang="ts">
    import MultiStepForm from "$lib/components/MultiStepForm.svelte"


</script>

{#snippet step1()}
    <div class="step">
        <h2>Transaction Details</h2>
        <p>Enter the details for your transaction.</p>
        <div class="remittance-fields">
        </div>
    </div>
{/snippet}

<main>
    <MultiStepForm
        steps={[step1]}
        method="POST"
    >

    </MultiStepForm>
</main>

<style>
    main {
        width: 95%;
        max-width: 900px;
        margin: 3.5rem auto;
        border: 1px solid saddlebrown;

        @media (max-width: 700px) {
            margin: 2.5rem auto;
        }
    }

    .step {
        overflow: hidden;
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
</style>