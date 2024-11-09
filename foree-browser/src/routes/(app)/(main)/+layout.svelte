<script lang="ts">
    const { children } = $props()
</script>


<nav class="header">
    Nulla ullamco labore sunt sunt officia in eu aliquip tempor consequat elit eiusmod Lorem officia.
    <p>aaaa</p>
</nav>

<div id="foree-main">
    <div class="container">
        <div class="nav-container">
            <a href="dashboard">Dashboard asdf asdfasdf asdfasdf asdfdfsa asdfasdf asdfasdf</a>
        </div>
        <div class="page-container">
            {@render children()}
        </div>
    </div>
</div>

<style>
    .header {
        z-index: 101;
        height: var(--foree-nav-height);
        padding: 0 var(--foree-page-padding-side);
        width: 100vw;
        margin: 0 auto;
        position:fixed;
        top:0;
        background-color: var(--foree-bg-1);

        &:after {
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

            &:after {
                background: linear-gradient(#0000, #0000000d);
                top: -4px;
                bottom: unset;
            }
        }
    }

    #foree-main {
        padding-top: var(--foree-banner-height);
        height: 100%;
        border: 1px solid salmon;
        /* margin: o auto;
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

</style>