<script lang="ts">
    import { page } from "$app/stores"
    import dropDownIcon from "$lib/assets/icons/drop_down_arrow.png"
    import { clickOutside } from "$lib/utils/use_directives"

    const { children } = $props()

    let desktopDropdownOn = $state(false)

    let mobileMenuOn = $state(false)
    
</script>

<nav class="header">
    <a class="desktop home-link" href="dashboard" title="Homepage" aria-label="Homepage"></a>
    <a class="mobile" href="dashboard" title="Homepage" aria-label="Homepage">X</a>
    <div class="desktop">
        <div>
            <div class="notifications">
                <a class:selected={$page.url.pathname === "/transaction"} href="transaction" title="Notifications" aria-label="notifications"></a>
            </div>
            <div 
                class="dropdown"
            >
                <button 
                    onclick={_ => {desktopDropdownOn = !desktopDropdownOn}} 
                    class:dropdown-content-on={desktopDropdownOn}
                    use:clickOutside={() => {
                        desktopDropdownOn = false
                    }}
                >
                    <p>XXXX XXXXXXXXXX xxx</p> <img src={dropDownIcon} alt=""/>
                </button>
                <nav class="dropdown-content" class:dropdown-content-on={desktopDropdownOn}>
                    <ul>
                        <li>
                            <a class="desktop home-link" href="dashboard" title="Homepage" aria-label="Homepage">Home</a>
                        </li>
                        <li>
                            <a href="transaction" class:selected={$page.url.pathname === "/transaction"} >Transaction</a>
                        </li>
                        <li>
                            <a class="desktop home-link" href="dashboard" title="Homepage" aria-label="Homepage">Sign Out</a>
                        </li>
                    </ul>
                </nav>
            </div>
        </div>
    </div>
    <div class="mobile">
        <div class="mobile-menu">
            <button aria-label="transfer"></button>
            <button aria-label="menu" class:open={mobileMenuOn} onclick={_ => mobileMenuOn = !mobileMenuOn} ></button>
        </div>
    </div>
</nav>

{#if mobileMenuOn}
<div id="mobile-menu-item" class="mobile">
    <div class="modal-overlay" aria-hidden="true"></div>
    <div class="menu">
        <div class="mobile-main-menu">
            <div class="menu-background ready" style:height="100%">
                <p>Laborum sunt ea irure culpa. Ad elit culpa eu aliqua elit minim officia aute Lorem. Consequat irure voluptate adipisicing qui minim irure duis in elit cillum sint nisi labore. Eu nostrud aliqua aliquip ullamco. Laboris sit duis magna dolor reprehenderit anim minim proident.

Duis non aliquip nulla et occaecat voluptate dolore ullamco labore labore ea. Ad adipisicing irure culpa anim nostrud proident. Culpa do enim aliqua ad est velit in quis do. Mollit officia nulla labore consectetur commodo. Cillum proident excepteur velit et sit ullamco veniam est. Ipsum laboris sunt et ullamco. Labore labore mollit voluptate sunt ex amet laboris id minim consectetur in.

Consectetur proident do excepteur dolore cupidatat commodo culpa sunt. Occaecat consequat voluptate amet reprehenderit proident sunt nulla amet voluptate ea culpa. Et incididunt pariatur et deserunt tempor excepteur veniam. Dolor quis deserunt eiusmod id aute occaecat qui aute sint labore adipisicing et proident. Ad incididunt proident et fugiat occaecat ipsum reprehenderit deserunt labore sunt sint exercitation do.

Sint ullamco quis ex id elit. Velit minim non exercitation non nostrud. Aute eu nulla ad esse amet pariatur aute pariatur. Duis tempor irure magna deserunt deserunt nisi aute aliqua sit magna id aute nulla. Aute magna labore voluptate pariatur do do velit sit minim.

Sit ullamco et non qui. Dolor irure consectetur reprehenderit in. Elit ipsum proident officia non amet deserunt consequat. Do elit pariatur dolor cillum excepteur. Quis in enim eu proident quis. Veniam ex dolor amet deserunt sunt reprehenderit dolore quis dolore laboris.

Do nulla commodo aliqua sunt magna Lorem veniam et. Proident magna magna culpa ex non aute voluptate laboris in fugiat proident. Id eu id veniam exercitation laborum irure eiusmod magna incididunt. Excepteur nisi nulla eiusmod sint. Elit irure quis tempor cillum deserunt excepteur ad qui excepteur proident ex incididunt.

Ad fugiat eu id cillum ullamco fugiat sunt consequat eiusmod adipisicing do incididunt elit occaecat. Eiusmod amet eiusmod est culpa. Incididunt ullamco et ad aliqua. Minim consectetur reprehenderit anim qui consequat non aute enim est ex labore eu quis id. Ex sit ut non nostrud voluptate aliqua ad pariatur do nisi.

Exercitation incididunt est veniam ut magna excepteur dolor velit sunt. Anim anim quis est proident dolore sit fugiat culpa minim laborum. Reprehenderit ut ullamco proident ea laborum adipisicing non do id aliqua occaecat amet. Anim officia amet et duis. Lorem commodo ullamco tempor Lorem irure ad ad veniam qui. Cillum dolor culpa aliquip sint duis fugiat irure duis aliquip reprehenderit officia et.

Dolor laboris est fugiat ut qui. Consectetur fugiat laboris eu ea occaecat id dolore sint ex officia. Ex excepteur dolor ea occaecat non culpa ex labore magna duis aute esse. Reprehenderit sint reprehenderit tempor velit excepteur occaecat.

Cillum ipsum Lorem incididunt nostrud id commodo tempor nostrud. Est eu est sunt quis. Incididunt ipsum quis velit irure. Eiusmod enim esse consectetur do esse cupidatat adipisicing pariatur laboris id cupidatat minim.</p>
            </div>
        </div>
    </div>
</div>
{/if}

<div id="foree-main">
    <div class="container">
        <div class="nav-container">
            <ul>
                <li><a href="dashboard" class:selected={$page.url.pathname === "/dashboard"}>Dashboard</a></li>
                <li><a href="transaction" class:selected={$page.url.pathname === "/transaction"} >Transaction</a></li>
            </ul>
        </div>
        <div class="page-container">
            {@render children()}
        </div>
    </div>
</div>

<style>

    .desktop {
        display: none;
        @media (min-width: 832px) {
            display: unset;
        }
    }

    .mobile {
        @media (min-width: 832px) {
            display: none;
        }
    }

    .header {
        z-index: 101;
        height: var(--foree-nav-height);
        padding: 0 var(--foree-page-padding-side);
        width: 100vw;
        margin: 0 auto;
        position:fixed;
        top:0;
        background-color: var(--foree-bg-1);
        display: grid;
        grid-template-columns: auto 1fr;
        box-sizing: border-box;

        &::after {
            content: "";
            background: linear-gradient(#0000000d, #0000);
            width: 100%;
            height: 4px;
            position: absolute;
            bottom: -4px;
            left: 0;
        }

        @media (max-width: 832px) {
            & {
                top: unset;
                bottom: 0;
            }

            &::after {
                background: linear-gradient(#0000, #0000000d);
                top: -4px;
                bottom: unset;
            }
        }
    }

    .header > a.mobile {
        background-color: var(--emerald-800);
        text-decoration: none;
        width: 50px;
        height: 50px;
        border-radius: 50%;
        font-size: xx-large;
        font-weight: 600;
        color: var(--slate-100);
        text-align: center;
        line-height: 50px;
        margin: auto 0;
        position: relative;
        -moz-box-shadow: 0 0 1px 1px var(--emerald-600);
        -webkit-box-shadow: 0 0 1px 1px var(--emerald-600);
        box-shadow: 0 0 1px 1px var(--emerald-600);

        &::after {
            content: "";
            background: var(--rose-700);
            height: 15px;
            width: 15px;
            right: 0;
            top: 0;
            position: absolute;
            border-radius: 50%;
            -moz-box-shadow: 0 0 1px 1px var(--rose-600);
            -webkit-box-shadow: 0 0 1px 1px var(--rose-600);
            box-shadow: 0 0 1px 1px var(--rose-600);
        }
    }

    .header .mobile-menu {
        display: flex;
        gap: 0.5rem;
        justify-content: end;
        width: 100%;
        height: 100%;
        align-items: center;

        button {
            background: transparent;
            width: 2.5rem;
            height: 2.5rem;
            background-size: 1.3rem;
            background-repeat: no-repeat;
            background-position: center;
            border-radius: var(--foree-border-radius);
            border-width: var(--foree-raised-width);
            border-color: var(--foree-raised-color);

            &:hover {
                background-color: var(--foree-raised-hover-color);
            }
        }

        & button[aria-label="transfer"] {
            background-image: url("$lib/assets/icons/send.png");
        }

        & button[aria-label="menu"] {
            background-image: url("$lib/assets/icons/more_bar.png");

            &.open {
                background-image: url("$lib/assets/icons/x.png");
            }
        }
    }

    .header > div.desktop > div {
        display: flex;
        justify-content: end;
        height: 100%;
    }

    .header > div.desktop > div > .notifications {
        width: var(--foree-nav-height);

        & > a {
            display: block;
            width: 100%;
            height: 100%;
            filter: hue-rotate(5deg);
            background-size: 27px;
            background-repeat: no-repeat;
            background-image: url("$lib/assets/icons/bell.png");
            background-position: center;

            &:hover {
                background-color: var(--foree-bg-4);
            }

            &.selected {
                background-color: var(--foree-bg-4);
            }
        }
    }

    .header > div.desktop > div > .dropdown {
        display: flex;
        justify-content: end;
        height: 100%;

        & button {
            display: flex;
            gap: 0.5rem;
            background: transparent;
            border: none;
            align-items: center;
            padding: 0 0.5rem;
            width: 150px;

            & p {
                font-size: large;
                text-wrap: nowrap;
                text-overflow: ellipsis;
                overflow: hidden;
            }

            & img {
                height: 17px;
                width: 17px;
                /* flex: 0 0 17px; */
            }

            &:hover {
                background-color: var(--foree-bg-4);
            }
        }
    }

    .header .desktop .dropdown {
        position: relative;

        & > .dropdown-content {
            position: absolute;
            opacity: 0;
            pointer-events: none;
            background-color: var(--foree-bg-2);
            border: 1px solid var(--slate-200);
            border-top: none;
            filter: var(--foree-shadow);
            top: var(--foree-nav-height);
            width: 150px;
            padding: 0rem 0rem 0.5rem;
            border-bottom-left-radius: 5px;
            border-bottom-right-radius: 5px;

            &.dropdown-content-on {
                pointer-events: all;
                opacity: 1;
            }
        }

        & button.dropdown-content-on {
            background-color: var(--foree-bg-4);
        }

        & a {
            display: block;
            padding: 0.5rem 0.3rem;
            text-decoration: none;
            color: var(--slate-600);
            font-weight: 500;

            &:hover, &.selected {
                background-color: var(--foree-bg-4);
            }

        }

        & li:last-child a {
            color: var(--rose-600);
            
            &:hover {
                background-color: var(--rose-100);
            }
        }

    }

    .header > a.desktop {
        background-size: 100% 100%;
        background-image: url("$lib/assets/images/foree_remittance_small_logo.svg");
        width: 50px;

        @media (min-width: 823px) {
            background-image: url("$lib/assets/images/foree_remittance_logo.svg");
            width: 7.5rem;
        }
    }

    #foree-main {
        padding-top: var(--foree-banner-height);
        height: 100%;
        border: 1px solid red;
        min-height: 100vh;
        /* margin: 0 auto;
        position: relative; */

        @media (min-width: 832px) {
            & {
                padding-top: var(--foree-nav-height);
                padding-bottom: var(--foree-banner-height);
            }
        }

        .container {
            --sidebar-menu-width: 16rem;
            --sidebar-width: var(--sidebar-menu-width);
            /* display: flex;
            flex-direction: column; */

            .nav-container {
                background: var(--foree-bg-2);
                display: none;


                @media (min-width: 832px) {
                    width: var(--sidebar-width);
                    height: calc(100vh - var(--foree-nav-height) - var(--foree-banner-height));
                    position: fixed;
                    left: 0px;
                    top: var(--foree-nav-height);
                    display: block;     
                    overflow: hidden;

                    &::after {
                        content: "";
                        background: linear-gradient(90deg, rgba(0,0,0,0), rgba(0,0,0,0.03));
                        width: 3px;
                        height: 100%;
                        position: absolute;
                        top: 0px;
                        right: 0px;
                    }
                }
            }

            .page-container {
                padding: var(--foree-page-padding-top) var(--foree-page-padding-side) var(--foree-page-padding-bottom);
                /* min-width: 0px !important; */

                @media (min-width: 832px) {
                    & {
                        padding-left: calc(var(--sidebar-width) + var(--foree-page-padding-side));
                    }
                }

                @media (min-width: 1536px) {
                    & {
                        --foree-page-padding-side: 6rem;
                    }
                }
            }
        }
    }

    .nav-container {
        & a {
            display: block;
            position: relative;
            height: var(--foree-side-menu-item-height);
            padding-left: var(--foree-page-padding-side);
            line-height: var(--foree-side-menu-item-height);
            text-wrap: nowrap;
            text-overflow: ellipsis;
            text-decoration: none;
            font-size: large;
            font-weight: 600;
            color: var(--slate-600);

            &:hover {
                background-color: var(--foree-bg-4);
            }

            @media (min-width: 832px) {
                    
                &.selected {
                    color: white;
                    background-color: var(--emerald-700);
                }

                &.selected::after {
                    --size: 1rem;
                    content: "";
                    width: var(--size);
                    height: var(--size);
                    background-color: var(--foree-bg-1);
                    position: absolute;
                    top: calc(var(--foree-side-menu-item-height)*.5 - var(--size)*.5);
                    right: calc(-.5*var(--size));
                    rotate: 45deg;
                    z-index: 2;
                    /* box-shadow: rgba(0, 0, 0, 0.12) 0px 0px 3px; */
                }
            }
        }
    }

    #mobile-menu-item > .modal-overlay {
        z-index: 99;
        opacity: .7;
        pointer-events: auto;
        width: 100%;
        height: 100%;
        height: 100dvh;
        background: var(--foree-bg-1);
        position: fixed;
        inset: 0 auto atuo 0;
    }

    #mobile-menu-item > .menu {
        z-index: 100;
        position: fixed;
        left: 0;
        bottom: var(--foree-nav-height, 0);
        filter: drop-shadow(0px 0px 15px #0000001a);
        width: 100%;
        height: 70vh;
        overflow: hidden;
        border-radius: 1rem 1rem 0 0;

        & > .mobile-main-menu {
            height: 100%;
            contain: layout paint;

            & > .menu-background {
                background: var(--foree-bg-2);
                will-change: height;
                position: absolute;
                width: 99.5%;
                left: 0;
                bottom: 0;
                height: 0%;
                border-radius: 1rem 1rem 0 0;
                transition-property: none;
                transition: .4s ease-in;
            }

            & > .ready {
                transition-property: height;
            }
        }
    }
</style>