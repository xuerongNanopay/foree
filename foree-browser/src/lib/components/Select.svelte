<script lang="ts">
    import { clickOutside } from "$lib/utils/use_directives";
    import type { HTMLInputAttributes } from "svelte/elements"

    interface Props extends HTMLInputAttributes {
        label?: string
        options: { value: string, title: string}[]
    }

    let { id, name, label, options }: Props = $props()
    let openDropdown = $state(false)
    let inputWidth = $state()
    let value = $state()

    function toggleDropdown() {
        openDropdown = !openDropdown
        console.log(openDropdown)
    }
    $inspect(inputWidth)
</script>

<div id="select">
    <label for={id}>{label ?? name ?? id}</label>
    <input 
        id={id} 
        name={name} 
        readonly 
        onclick={toggleDropdown}
        value={value}
        placeholder="please select..."
        bind:offsetWidth={inputWidth}
        use:clickOutside={() => {
            if (openDropdown) openDropdown = false
        }}
    />

    {#if openDropdown}
        <div 
            class="dropdown-wrap"
            style:--dropdown-width={inputWidth+"px"}
        >
            <div class="dropdown">
                <ul>
                    {#each options as option}
                        <li><button type="button" onclick={_=>value=option.value}>{option.title}</button></li>
                    {/each}
                </ul>
            </div>
        </div>
    {/if}
</div>

<style>

    #select input {
        cursor: pointer;
    }

    #select .dropdown-wrap {
        position: fixed;
        margin-top: .5rem;
        width: var(--dropdown-width);
        & .dropdown {
            padding: .5rem 0;
            border: 1px solid var(--emerald-800);
            background-color: #fff;
            border-radius: 4px;
            width: 100%;
            z-index: 99;
            position: absolute;
        }
    }

    #select .dropdown button {
        padding: .5rem .2rem;
        background: none;
        border: none;
        display: block;
        width: 100%;
        text-align: left;
        &:hover{
            background: var(--slate-200);
        }
    }

</style>