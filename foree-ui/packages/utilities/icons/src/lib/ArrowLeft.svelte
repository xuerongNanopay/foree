<script lang="ts">
  import type { SVGAttributes } from 'svelte/elements';
  let { strokeWidth = 1.5, fill="currentColor", ...props }: SVGAttributes<EventTarget> = $props();
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
    d="M9.57 5.92993L3.5 11.9999L9.57 18.0699"
    stroke="currentColor"
    strokeLinecap="round"
    strokeLinejoin="round"
    strokeMiterlimit="10"
    strokeWidth={strokeWidth}
  />
  <path
    d="M20.5 12H3.67004"
    stroke="currentColor"
    strokeLinecap="round"
    strokeLinejoin="round"
    strokeMiterlimit="10"
    strokeWidth={strokeWidth}
  />
</svg>