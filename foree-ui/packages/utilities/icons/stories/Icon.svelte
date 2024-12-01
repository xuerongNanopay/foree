<script lang="ts">
  import './icon.css';
  import * as allIcons from '../src/lib';
  import type { Component } from 'svelte';
  import type { SVGAttributes } from 'svelte/elements';

  const icons: {name: string, component: Component<SVGAttributes<EventTarget>>}[] = []
  
  for (let n of Object.keys(allIcons)) {
    for (let _ of Array(140)) {
      icons.push({name: n, component: allIcons[n]})
    }
  }
  const { size= 'medium' } = $props();

  function formatCellName(name: string): string {
    return (name.slice(0,1).toLocaleLowerCase() + name.slice(1)).replace('Icon', '').replace(/[A-Z]/g, letter => `_${letter.toLowerCase()}`)
  }
</script>

{#snippet cell(icon)}
  <div class="cell">
    <div class="cell__icon__container">
      <icon.component/>
    </div>
    <small class="cell__desp">{formatCellName(icon.name)}</small>
  </div>
{/snippet}

<div>
  {#each icons as icon}
    <!-- <div class="cell">
      <div class="cell__icon__container">
        <icon.component/>
      </div>
      <small class="cell__desp">{formatCellName(icon.name)}</small>
    </div> -->
    {@render cell(icon)}
  {/each}
</div>