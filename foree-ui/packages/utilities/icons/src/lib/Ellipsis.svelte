<script lang="ts">
  import type { SVGAttributes } from 'svelte/elements';
  let { color="currentColor", ...props }: SVGAttributes<EventTarget> = $props();
</script>

<svg
  aria-hidden="true"
  fill="none"
  height="1em"
  viewBox="0 0 24 24"
  width="1em"
  {...props}
>
<circle cx="12" cy="12" fill={color} r="1" />
<circle cx="19" cy="12" fill={color} r="1" />
<circle cx="5" cy="12" fill={color} r="1" />
</svg>