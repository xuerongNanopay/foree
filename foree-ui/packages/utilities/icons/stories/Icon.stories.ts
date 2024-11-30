import type { Meta, StoryObj } from '@storybook/svelte';
import Icon from './Icon.svelte';

const meta: Meta<typeof Icon> = {
    title: 'icons/Button',
    component: Icon,
};

type Story = StoryObj<typeof meta>;

export default meta;

export const Primary: Story = {
    args: {
    },
  };