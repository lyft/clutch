import { create } from '@storybook/theming/create';

export default create({
  base: 'light',

  colorPrimary: '#ffffff',
  colorSecondary: '#02acbe',

  // UI
  appBg: '#ffffff',
  appContentBg: '#ffffff',
  appBorderColor: '#02acbe',
  appBorderRadius: 4,

  // Typography
  fontBase: '"Open Sans", sans-serif',
  fontCode: 'monospace',

  // Text colors
  textColor: '#2D3F50',

  // Toolbar default and active colors
  barTextColor: '#2D3F50',
  barSelectedColor: '#02acbe',
  barBg: '#D7DADB',

  brandTitle: 'Clutch Storybook',
  brandUrl: 'https://clutch.sh',
  brandImage: 'https://clutch.sh/img/navigation/logo.svg',
});