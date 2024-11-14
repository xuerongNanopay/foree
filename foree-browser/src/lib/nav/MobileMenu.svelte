<script lang="ts">
    import { expoOut, quintOut } from 'svelte/easing'
    import type { TransitionConfig } from 'svelte/transition'

    import ModalOverlay from "$lib/components/ModalOverlay.svelte"
    import { reduced_motion } from "$lib/stores/reduced_motion"
    import { trap } from "$lib/actions/focus";
    import MobileSubMenu from "./MobileSubMenu.svelte"
	import { mobileNavigation } from './MobileNavigation'

    interface Props {
		// links: NavigationLink[];
		// current: NavigationLink | undefined;
		onclose: () => void;
	}

	// let { links, current, onclose }: Props = $props();
    let { onclose }: Props = $props()


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
		></div>
        <div class="clip-wapper">
            <div class="viewport">
                <MobileSubMenu mobileNavigation={mobileNavigation}></MobileSubMenu>
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
		will-change: height;
	}

    .clip-wapper {
		width: 100%;
		height: 100%;
		transition: clip-path 0.3s cubic-bezier(0.23, 1, 0.32, 1);
        /* Why this affect background absolute? */
		will-change: clip-path;
        overflow-x: hidden;
		overflow-y: scroll;
	}

    .viewport {
		margin: 1rem 0;
	}
</style>