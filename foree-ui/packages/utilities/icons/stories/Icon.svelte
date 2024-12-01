<script lang="ts">
  import './icon.css';
  import * as allIcons from '../src/lib';
  import type { Component } from 'svelte';
  import type { SVGAttributes } from 'svelte/elements';

  const icons: {name: string, component: Component<SVGAttributes<EventTarget>>}[] = []

  type Prop = {
    size: 'small' | 'medium' | 'large' | 'xlarge',
    color?: 'string',
  }
  
  for (let n of Object.keys(allIcons)) {
    icons.push({name: n, component: allIcons[n]})
  }
  const { size='medium', color }: Prop = $props();

  function formatCellName(name: string): string {
    return (name.slice(0,1).toLocaleLowerCase() + name.slice(1)).replace('Icon', '').replace(/[A-Z]/g, letter => `-${letter.toLowerCase()}`)
  }
</script>

<!-- TODO: why typescript error -->
<!-- {#snippet cell(icon)}
  <div class="cell">
    <div class="cell__icon__container">
      <icon.component class={`cell__icon--${size}`} fill={fill}/>
    </div>
    <small class="cell__desp">{formatCellName(icon.name)}</small>
  </div>
{/snippet} -->

<div>
  {#each icons as icon}
    <!-- {@render cell(icon)} -->
    <div class="cell">
      <div class="cell__icon__container">
        <icon.component class={`cell__icon--${size}`} color={color}/>
      </div>
      <small class="cell__desp">{formatCellName(icon.name)}</small>
    </div>
  {/each}
</div>