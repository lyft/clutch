import type { ThemeConfig as AlgoliaThemeConfig } from '@docusaurus/theme-search-algolia';
import type { ThemeConfig as BaseThemeConfig } from '@docusaurus/preset-classic';

export interface FooterLink {
  title: string;
  items: {
    label: string;
    to: string;
  }[];
}

export interface ThemeConfig extends BaseThemeConfig {
  image: string;
  colorMode: {
    disableSwitch?: boolean;
  };
  algolia: AlgoliaThemeConfig;
  prism: {
    additionalLanguages: string[];
  };
  navbar: {
    logo: {
      alt: string;
      src: string;
    };
    hideOnScroll?: boolean;
  };
  footer: {
    style: 'light' | 'dark';
    logo: {
      src: string;
    };
    links: FooterLink[];
    copyright: string;
  }
}