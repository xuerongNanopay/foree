<script lang="ts">
  import type { SVGAttributes } from 'svelte/elements';
  let { color="currentColor", ...props }: SVGAttributes<EventTarget> = $props();
</script>

<svg
  aria-hidden="true"
  fill="none"
  focusable="false"
  height="1em"
  role="presentation"
  viewBox="0 0 24 24"
  width="1em"
  {...props}
>
  <path
    d="M12 22C17.5228 22 22 17.5228 22 12C22 6.47715 17.5228 2 12 2C6.47715 2 2 6.47715 2 12C2 17.5228 6.47715 22 12 22Z"
    stroke={color}
  />
  <path
    d="M8.46997 13.26L12 9.73999L15.53 13.26"
    stroke={color}
  />
</svg>