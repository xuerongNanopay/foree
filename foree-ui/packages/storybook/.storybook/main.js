import { join, dirname } from 'path';

/**
 * This function is used to resolve the absolute path of a package.
 * It is needed in projects that use Yarn PnP or are set up within a monorepo.
 */
function getAbsolutePath(value) {
  return dirname(require.resolve(join(value, 'package.json')));
}

/** @type { import('@storybook/sveltekit').StorybookConfig } */
const config = {
  stories: ['../src/**/*.mdx', '../src/**/*.stories.@(js|ts|svelte)'],
  addons: [
    // turborepo + pnpm won't work with storybook8
    // getAbsolutePath('@storybook/addon-svelte-csf'),
    // getAbsolutePath('@storybook/addon-essentials'),
    // getAbsolutePath('@chromatic-com/storybook'),
    // getAbsolutePath('@storybook/addon-interactions'),
    '@storybook/addon-svelte-csf',
    '@storybook/addon-essentials',
    '@chromatic-com/storybook',
    '@storybook/addon-interactions',
  ],
  framework: {
    // name: getAbsolutePath('@storybook/sveltekit'),
    name: '@storybook/sveltekit',
    options: {}
  }
};
export default config;
