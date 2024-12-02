<script lang="ts">
  import type { SVGAttributes } from 'svelte/elements';
  interface MoonProp extends SVGAttributes<EventTarget>  {
    filled?: boolean;
  }
  let { filled, color="currentColor", ...props }: MoonProp = $props();
</script>

{#if filled}
<svg
  aria-hidden="true"
  focusable="false"
  height="1em"
  role="presentation"
  viewBox="0 0 24 24"
  width="1em"
  {...props}
>
  <path
    d="M21.53 15.93c-.16-.27-.61-.69-1.73-.49a8.46 8.46 0 01-1.88.13 8.409 8.409 0 01-5.91-2.82 8.068 8.068 0 01-1.44-8.66c.44-1.01.13-1.54-.09-1.76s-.77-.55-1.83-.11a10.318 10.318 0 00-6.32 10.21 10.475 10.475 0 007.04 8.99 10 10 0 002.89.55c.16.01.32.02.48.02a10.5 10.5 0 008.47-4.27c.67-.93.49-1.519.32-1.79z"
    fill={color}
  />
</svg>
{:else}
<svg
  aria-hidden="true"
  focusable="false"
  height="1em"
  role="presentation"
  viewBox="0 0 512 512"
  width="1em"
  {...props}
>
  <path
    d="M160 136c0-30.62 4.51-61.61 16-88C99.57 81.27 48 159.32 48 248c0 119.29 96.71 216 216 216 88.68 0 166.73-51.57 200-128-26.39 11.49-57.38 16-88 16-119.29 0-216-96.71-216-216z"
    fill="none"
    stroke={color}
  />
</svg>
{/if}