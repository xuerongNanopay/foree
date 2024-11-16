<script lang="ts">
    import type { Snippet } from "svelte"
    import type { HTMLFormAttributes } from 'svelte/elements'

    interface Props extends HTMLFormAttributes {
        steps: Snippet<[]>[];
    }

    let { steps, ...rest }: Props = $props()
    let cur = $state(0)

    function goNext() {
        cur = cur < steps.length-1 ? cur + 1 : cur
    }

    function goPrev() {
        cur = cur > 0 ? cur - 1 : cur
    }

    let stepsDiv: HTMLDivElement;

    $effect(
        () => {
            const stepOffsetWidth = stepsDiv.offsetWidth/steps.length
            stepsDiv.style.marginLeft = (cur*stepOffsetWidth*-1)+"px"
        }
    )
</script>

<form {...rest}>
    <div style:width={steps.length*100 + "%"} bind:this={stepsDiv}>
        {#each steps as step}
            <div>
                {@render step()}
                <div class="btns">
                    {#if cur > 0}
                        <button class="pre" type="button" onclick={goPrev}>Prev</button>
                    {/if}
                    {#if cur < steps.length-1}
                        <button type="button" onclick={goNext}>Next</button>
                    {/if}
                    {#if cur === steps.length-1}
                        <button >Submit</button>
                    {/if}
                </div>
            </div>
        {/each}
    </div>
</form>


<style>
    form, form * {
        box-sizing: border-box;
    }

    form {
        overflow: hidden;

        & > div:first-child {
            display: flex;
            transition: margin-left 200ms ease-in-out;
            & > div {
                width: 100%;
            } 
        }
    }

    .btns {
        margin-top: 3rem;
        display: flex;
        gap: 0.5rem;

        @media (max-width: 700px) {
            & {
                flex-direction: column;
            }
        }

        & button[type="button"] {
            display: block;
            width: 100%;
            background-color: var(--primary-color);
            border: 0px;
            padding: 0.75rem 0;
            border-radius: 0.25rem;
            color: white;
            font-size: 1em;
            font-weight: 600;

            &.pre {
                background-color: transparent;
                color:  var(--primary-color);
                border: 2px solid  var(--primary-color);
            }

            &:disabled {
                opacity: 0.6;
            }
        }
    }
</style>