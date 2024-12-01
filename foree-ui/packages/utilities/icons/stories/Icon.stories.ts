import type { Meta, StoryObj } from '@storybook/svelte';
import Icon from './Icon.svelte';

const meta: Meta<typeof Icon> = {
  title: 'icons',
  component: Icon,
};

type Story = StoryObj<typeof Icon>;

export const SupportIcons: Story = {
  argTypes: {
    size: {
      control: { type: 'select' },
      options: ['small', 'medium', 'large', 'xlarge'],
    },
    color: {
      control: { type: 'color' },
    },
    filled: {
      control: { type: 'boolean' },
    }
  },
};

export default meta;