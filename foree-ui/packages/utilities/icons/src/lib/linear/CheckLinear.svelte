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
  stroke={color}
  viewBox="0 0 24 24"
  width="1em"
  {...props}
>
  <polyline points="20 6 9 17 4 12" />
</svg>