<script lang="ts">
    import { page } from "$app/stores"
    import UpArrow from "$lib/assets/icons/up_arrow.png"
    import type { MobileNavigations } from "$lib/types"

    interface Props {
		mobileNavigations: MobileNavigations
	}

    let { mobileNavigations }: Props = $props()
</script>

<ul class="accordion">
    {#each mobileNavigations as subMenu, id}
        <li>
            <input type="checkbox" name={""+id} id={""+id} checked={!!subMenu.defaultActive}>
            <label for={""+id}>
                <p>{subMenu.subMenuTitle}</p>
                <img src={UpArrow} alt="">
            </label>
            <ul class="content">
                {#each subMenu.navigations as navLik}
                    <li><a href={navLik.href} class:selected={$page.url.pathname === navLik.href}>
                        <img src={navLik.icon} alt=" "/>
                        <p>{navLik.title}</p>
                    </a></li>
                {/each}
            </ul>
        </li>
    {/each}
</ul>

<style>
    .accordion {
        padding: 0 1rem;
    }

    .accordion label {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: .5rem 0 .3rem;
        border-bottom: 1px solid var(--slate-200);
        &:hover {
            background-color: var(--foree-bg-4);
        }

        & p {
            font-size: large;
            font-weight: 500;
        }
        & img {
            height: 17px;
            width: 17px;
            transition: all .5s ease-in-out;
        }
    }

    .accordion input[type="checkbox"] {
        display: none;
    }

    .accordion .content {
        max-height: 0;
        overflow: hidden;
        transition: max-height .5s ease-in-out;

        & a {
            display: flex;
            align-items: center;
            gap: .5rem;
            text-decoration: none;
            padding: .5rem 1rem;
            border-radius: 7px;
            &:hover {
                background-color: var(--foree-bg-4);
            }

            &.selected {
                background-color: var(--foree-bg-4);
            }

            & img {
                width: 14px;
                height: 14px;
            }
        }
    }

    .accordion input[type="checkbox"]:checked {
        & ~ .content {
            max-height: 300px;
        }

        & ~ label > img {
            transform: rotate(180deg);
        }
    }
</style>