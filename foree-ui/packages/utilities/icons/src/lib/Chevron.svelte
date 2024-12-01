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
  d="M15.5 19l-7-7 7-7"
  stroke={color}
/>
</svg>