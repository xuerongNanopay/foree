<script lang="ts">
  import './icon.css';
  import * as allIcons from '../src/lib';
  import type { Component } from 'svelte';
  import type { SVGAttributes } from 'svelte/elements';

  const icons: {name: string, component: Component<SVGAttributes<EventTarget>>}[] = []
  
  for (let n of Object.keys(allIcons)) {
    icons.push({name: n, component: allIcons[n]})
  }
  const { size= 'medium' } = $props();

</script>

{#each icons as icon}
  <div class="cell">
    <div class="cell__icon__container">
      <icon.component/>
    </div>
    <!-- <small class="cell__desp">{icon.name.replace('Icon', '').toLowerCase()}</small> -->
    <small class="cell__desp">{'1111111111111'}</small>
  </div>
{/each}