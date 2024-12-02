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
  <path d="M0 0h24v24H0z" fill="none" stroke="none" />
  <path d="M8 9l4 -4l4 4" />
  <path d="M16 15l-4 4l-4 -4" />
</svg>