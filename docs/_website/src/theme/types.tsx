import type { ThemeConfig as AlgoliaThemeConfig } from "@docusaurus/theme-search-algolia";
import type { ThemeConfig as BaseThemeConfig } from "@docusaurus/preset-classic";
import type { DocusaurusConfig } from "@docusaurus/types";

interface HeroButtonConfig {
  url: string;
  text: string;
}

export interface FeatureConfig {
  title: string;
  imageUrl: string;
  description?: string;
}

export interface FeaturesConfig {
  title: string;
  featureList: FeatureConfig[];
}

export interface DemoConfig {
  lines: string[];
  cta: {
    text: string;
    link: string;
  };
}

export interface ConsolidationConfig {
  snippets: string[];
}

export interface HeroConfig {
  description: string;
  buttons: {
    first: HeroButtonConfig;
    second: HeroButtonConfig;
  };
}

export interface SiteConfig extends DocusaurusConfig {
  customFields: {
    sections: {
      features: FeaturesConfig;
      demo: DemoConfig;
      consolidation: ConsolidationConfig;
    };
    tagDescription: string;
    hero: HeroConfig;
    archivalNotice?: {
      enabled: boolean;
      title: string;
      message: string;
    };
  };
}

export interface FooterLinkItem {
  label: string;
  to: string;
  html?: string;
  href?: string;
  prependBaseUrlToHref?: boolean;
}

export interface FooterLink {
  title: string;
  items: Array<{
    label: string;
    to: string;
    html?: string;
  }>;
}

export interface FooterConfig {
  style: "light" | "dark";
  logo: {
    src: string;
  };
  links?: FooterLink[];
  copyright: string;
}

export interface ThemeConfig extends BaseThemeConfig {
  image?: string;
  colorMode?: {
    disableSwitch?: boolean;
  };
  algolia?: AlgoliaThemeConfig;
  prism?: {
    additionalLanguages: string[];
  };
  navbar?: {
    logo: {
      alt: string;
      src: string;
    };
    hideOnScroll?: boolean;
  };
  footer?: FooterConfig;
}
