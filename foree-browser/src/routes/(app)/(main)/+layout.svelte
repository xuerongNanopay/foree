<script lang="ts">
    const { children } = $props()
</script>


<nav class="header">
    Header
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

        @media (max-width: 831px) {
            & {
                top: unset;
                bottom: 0;
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
                background: red;
                display: none;


                @media (min-width: 832px) {
                    width: var(--sidebar-width);
                    height: calc(100vh - var(--foree-nav-height) - var(--foree-banner-height));
                    border: 1px solid rebeccapurple;
                    position: fixed;
                    left: 0px;
                    top: var(--foree-nav-height);
                    display: block;     
                    overflow: hidden;   
                }
            }

            .page-container {
                padding: var(--foree-page-padding-top) var(--foree-page-padding-side) var(--foree-page-padding-bottom);
            }
        }
    }

</style>