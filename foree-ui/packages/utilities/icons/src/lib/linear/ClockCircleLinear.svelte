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
<g fill="none" stroke={color}>
  <circle cx="12" cy="12" r="10" />
  <path d="M12 8v4l2.5 2.5"/>
</g>
</svg>