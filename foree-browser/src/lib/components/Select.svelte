<script lang="ts">
    import type { HTMLInputAttributes } from "svelte/elements"

    interface Props extends HTMLInputAttributes {
        label?: string
    }

    let { id, name, label }: Props = $props()
    let openDropdown = $state(true)

    function toggleDropdown() {
        openDropdown = !openDropdown
        console.log(openDropdown)
    }
</script>

<div id="select">
    <label for={id}>{label ?? name ?? id}</label>
    <input id={id} name={name} readonly onclick={toggleDropdown} placeholder="please select..."/>

    {#if openDropdown}
        <div class="dropdown-wrap">
            <div class="dropdown">
                <ul>
                    <li><p>aaaaaa</p></li>
                    <li><p>aaaaaa</p></li>
                    <li><p>aaaaaa</p></li>
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
        position: relative;
        & .dropdown {
            margin-top: .5rem;
            border: 1px solid var(--emerald-800);
            background-color: #fff;
            border-radius: 7px;
            width: 100%;
            z-index: 99;
            position: absolute;
        }
    }

</style>