<script lang="ts">
    import { page } from '$app/stores'
    const { children } = $props()
</script>


<nav class="header">
    <a class="desktop home-link" href="dashboard" title="Homepage" aria-label="Homepage"></a>
    <a class="mobile" href="dashboard" title="Homepage" aria-label="Homepage">X</a>
    <div class="desktop"></div>
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

    .header {
        & a[class*="mobile"] {
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

    .header .home-link {
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
                    background-color: var(--foree-bg-4);
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