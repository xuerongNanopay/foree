<script lang="ts">
  import type { SVGAttributes } from 'svelte/elements';
  interface MailProp extends SVGAttributes<EventTarget>  {
    filled?: boolean;
  }
  let { filled, color="currentColor", ...props }: MailProp = $props();
</script>

{#if filled}
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
    d="M17 3.5H7C4 3.5 2 5 2 8.5V15.5C2 19 4 20.5 7 20.5H17C20 20.5 22 19 22 15.5V8.5C22 5 20 3.5 17 3.5ZM17.47 9.59L14.34 12.09C13.68 12.62 12.84 12.88 12 12.88C11.16 12.88 10.31 12.62 9.66 12.09L6.53 9.59C6.21 9.33 6.16 8.85 6.41 8.53C6.67 8.21 7.14 8.15 7.46 8.41L10.59 10.91C11.35 11.52 12.64 11.52 13.4 10.91L16.53 8.41C16.85 8.15 17.33 8.2 17.58 8.53C17.84 8.85 17.79 9.33 17.47 9.59Z"
    fill={color}
  />
</svg>
{:else}
<svg
  aria-hidden="true"
  focusable="false"
  height="1em"
  role="presentation"
  viewBox="0 0 24 24"
  width="1em"
  {...props}
>
  <g
    fill="none"
    stroke={color}
  >
    <path d="M12 20.5H7c-3 0-5-1.5-5-5v-7c0-3.5 2-5 5-5h10c3 0 5 1.5 5 5v3" />
    <path d="M17 9l-3.13 2.5a3.166 3.166 0 01-3.75 0L7 9M19.21 14.77l-3.539 3.54a1.232 1.232 0 00-.3.59l-.19 1.35a.635.635 0 00.76.76l1.35-.19a1.189 1.189 0 00.59-.3l3.54-3.54a1.365 1.365 0 000-2.22 1.361 1.361 0 00-2.211.01zM18.7 15.28a3.185 3.185 0 002.22 2.22" />
  </g>
</svg>
{/if}