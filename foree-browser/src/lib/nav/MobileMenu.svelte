<script lang="ts">
    import ModalOverlay from "$lib/components/ModalOverlay.svelte";
    import type { NavigationLink } from "$lib/types"
    import type { TransitionConfig } from 'svelte/transition'

    interface Props {
		// links: NavigationLink[];
		// current: NavigationLink | undefined;
		onclose: () => void;
	}

	// let { links, current, onclose }: Props = $props();

    let { onclose }: Props = $props()

</script>

<ModalOverlay {onclose} />

<!-- <div id="mobile-menu-item" class="mobile">
    <div class="modal-overlay" aria-hidden="true"></div>
    <div class="menu">
        <div class="mobile-main-menu">
            <div class="menu-background ready" style:height="100%"></div>
            <div class="clip">
            </div>
        </div>
    </div>
</div> -->