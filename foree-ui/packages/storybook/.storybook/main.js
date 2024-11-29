/** @type { import('@storybook/svelte-vite').StorybookConfig } */
const config = {
  stories: [
    "../../utilities/icons/stories/**/*.stories.@(js|jsx|ts|tsx)",
  ],
  framework: '@storybook/svelte-vite',
  addons: [
    "@storybook/addon-svelte-csf",
    "@storybook/addon-essentials",
    "@chromatic-com/storybook",
    "@storybook/addon-interactions"
  ],
};
 
export default config;