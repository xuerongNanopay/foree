<script lang="ts">
    import { page } from "$app/stores"
    import dropDownIcon from "$lib/assets/icons/drop_down_arrow.png"
    import { clickOutside } from "$lib/utils/use_directives"

    const { children } = $props()

    let dropdownOn = $state(false)
    
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
                    onclick={_ => {dropdownOn = !dropdownOn}} 
                    class:dropdown-content-on={dropdownOn}
                    use:clickOutside={() => {
                        dropdownOn = false
                    }}
                >
                    <p>XXXX XXXXXXXXXX xxx</p> <img src={dropDownIcon} alt=""/>
                </button>
                <nav class="dropdown-content" class:dropdown-content-on={dropdownOn}>
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
    <div class="mobile"></div>
</nav>

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

    .header > a[class*="mobile"] {
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

    .header > div[class*="desktop"] > div {
        display: flex;
        justify-content: end;
        height: 100%;
    }

    .header > div[class*="desktop"] > div > .notifications {
        width: var(--foree-nav-height);

        & > a {
            display: block;
            width: 100%;
            height: 100%;
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

    .header > div[class*="desktop"] > div > .dropdown {
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

    .header > a[class*="desktop"]  {
        background-size: 100% 100%;
        background-image: url("$lib/assets/images/foree_remittance_small_logo.svg");
        width: 50px;

        @media (min-width: 823px) {
            background-image: url("$lib/assets/images/foree_remittance_logo.svg");
            width: 7.5rem;
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
</style>