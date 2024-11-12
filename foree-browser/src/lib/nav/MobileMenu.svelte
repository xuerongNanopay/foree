<script lang="ts">
    import type { NavigationLink } from "$lib/types"
    import { expoOut, quintOut } from 'svelte/easing'
    import type { TransitionConfig } from 'svelte/transition'

    import ModalOverlay from "$lib/components/ModalOverlay.svelte"
    import { reduced_motion } from "$lib/stores/reduced_motion"
    import { trap } from "$lib/actions/focus";

    interface Props {
		// links: NavigationLink[];
		// current: NavigationLink | undefined;
		onclose: () => void;
	}

	// let { links, current, onclose }: Props = $props();
    let current = $state.raw<NavigationLink | undefined>();
    let { onclose }: Props = $props()
    let universal_menu_inner_height = $state(0);
    let menu_height = $state(0);
    let show_context_menu = $state(false);
    let ready = $state(false);
    let universal_menu: HTMLElement | undefined = $state();


    function popup(node: HTMLElement, { duration = 400, easing = expoOut } = {}): TransitionConfig {
		const height = node.clientHeight;

		return {
			css: (t, u) =>
				$reduced_motion
					? `opacity: ${t}`
					: `transform: translate3d(0, ${(height * u) / 0.9}px, 0) scale(${0.9 + 0.1 * t})`,
			easing,
			duration
		};
	}

</script>

<ModalOverlay {onclose} />

<div class="mobile-menu">
    <div class="mobile-menu-content" transition:popup={{ duration: 400, easing: quintOut }}>
		<div
			class="mobile-menu-content-background"
			class:ready
			style:height='100%'
		></div>
        <div 
            class="clip-wapper"
            style:--height-difference="{menu_height - universal_menu_inner_height}px"
        >
            <div
                class="viewport"
                class:reduced-motion={$reduced_motion}
                class:offset={show_context_menu}
                bind:clientHeight={menu_height}
                onscroll={(e:Event) => {e.stopPropagation()}}
            >
                <div class="universal" inert={false} bind:this={universal_menu}>
                    <div class="contents" bind:clientHeight={universal_menu_inner_height} onscroll={(e:Event) => {e.stopPropagation()}}>
						<ul>
							<li><a href="/chat">Discord</a></li>
							<li><a href="https://bsky.app/profile/sveltesociety.dev">Bluesky</a></li>
							<li><a href="https://github.com/sveltejs/svelte">GitHub</a></li>
						</ul>
                    </div>
                </div>
            </div>
        </div>
    </div>
</div>

<style>
    .mobile-menu {
		display: block;
		position: fixed;
		left: 0px;
		bottom: var(--bottom, var(--foree-nav-height));
		z-index: 100;
		width: 100%;
		height: 70vh;
		border-radius: 1rem 1rem 0 0;
		overflow-y: hidden;
		overflow-x: hidden;
		/* pointer-events: none; */
		transform: translate3d(0, 0, 0);
		filter: var(--foree-shadow);
	}

    .mobile-menu-content {
		height: 100%;
		contain: layout paint;
		transform: translateZ(0);
		backface-visibility: hidden;
	}

    .mobile-menu-content-background {
		position: absolute;
		width: 100%;
		left: 0;
		bottom: 0;
		height: 99.5%;
		border-radius: 1rem 1rem 0 0;
		background: var(--background, var(--foree-bg-2));
        background-color: red;
		will-change: height;
		/* transition: 0.3s var(--quint-out);
		transition-property: none; */

		&.ready {
			/* transition-property: height; */
		}

		/* :root.dark & {
			border-top: solid 1px var(--sk-raised-highlight);
		} */
	}

    .clip-wapper {
        background: seagreen;
		width: 100%;
		height: 100%;
		transition: clip-path 0.3s cubic-bezier(0.23, 1, 0.32, 1);
        /* Why this affect background absolute. */
		will-change: clip-path;
	}

    .viewport {
		position: relative;
		bottom: -1px;

		/* display: grid; */
		width: 200%;
		height: 100%;
		grid-template-columns: 50% 50%;
		transition: transform 0.3s cubic-bezier(0.23, 1, 0.32, 1);
		grid-auto-rows: 100%;

		&.reduced-motion {
			transition-duration: 0.01ms;
		}

		&.offset {
			transform: translate3d(-50%, 0, 0);
		}

		& > * {
			overflow-y: scroll;
			transition: inherit;
			transition-property: transform, opacity;
		}

        & * {
            /* overscroll-behavior: contain; */
        }

		/* & :global(a) {
			position: relative;
			padding: 0.3rem 0;
			color: inherit;
			font: var(--sk-font-ui-medium);
			width: 100%;
			height: 100%;
		} */
	}

    .universal .contents {
		position: absolute;
		width: 50%;
		bottom: 0;
		padding: 1rem var(--foree-page-padding-side);
		max-height: 70vh;
		overflow-y: scroll;

		button {
			/* width: 2.6rem; */
			height: 2.6rem;
		}
	}
</style>