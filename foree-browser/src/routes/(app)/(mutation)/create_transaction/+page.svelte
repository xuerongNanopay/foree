<script lang="ts">
    import MultiStepForm from "$lib/components/MultiStepForm.svelte"
    import Select from "$lib/components/Select.svelte";


</script>

{#snippet step1()}
    <div class="step">
        <h2>Transaction Details</h2>
        <p>Enter the details for your transaction.</p>
        <div class="remittance-fields">
            <Select id="from" name="from" options={[{value: "33", title: "AAA"}]}></Select>
        </div>
    </div>
{/snippet}

<!-- {#snippet step2()}
    <div class="step">
        <h2>Transaction Details</h2>
        <p>Enter the details for your transaction.</p>
        <div class="remittance-fields">
            <Select id="from" name="from"></Select>
        </div>
    </div>
{/snippet} -->

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
</style>