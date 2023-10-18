import React from "react";
import clsx from "clsx";
import Link from "@docusaurus/Link";
import useDocusaurusContext from "@docusaurus/useDocusaurusContext";
import useBaseUrl from "@docusaurus/useBaseUrl";
import type { NavbarLogo } from "@docusaurus/theme-common";
import { useThemeConfig } from "@docusaurus/theme-common";

import styles from "./styles.module.css";
import {
  FooterLink,
  FooterLinkItem as IFoorterLinkitem,
  ThemeConfig,
} from "../types";

interface SocialLink {
  icon: string;
  href: string;
}

const socialLinks = [
  {
    icon: "fe fe-github",
    href: "https://github.com/lyft/clutch",
  },
  {
    icon: "fe fe-slack",
    href: "https://join.slack.com/t/lyftoss/shared_invite/zt-casz6lz4-G7gOx1OhHfeMsZKFe1emSA",
  },
] as SocialLink[];

function FooterLinkItem({
  to,
  href,
  label,
  prependBaseUrlToHref = false,
  ...props
}: IFoorterLinkitem): JSX.Element {
  const toUrl = useBaseUrl(to);
  const normalizedHref = useBaseUrl(href);

  return (
    <Link
      className="footer__link-item"
      {...(href !== undefined
        ? {
            target: "_blank",
            rel: "noopener noreferrer",
            href: prependBaseUrlToHref ? normalizedHref : href,
          }
        : {
            to: toUrl,
          })}
      {...props}
    >
      {label}
    </Link>
  );
}

function Logo({ ...props }): JSX.Element {
  const {
    navbar: { logo = {} as NavbarLogo },
  } = useThemeConfig();
  const lyftLogoUrl = useBaseUrl("img/navigation/lyft-logo.svg");
  const logoImageUrl = useBaseUrl("img/navigation/logo.svg");

  return (
    <div className={clsx("navbar__brand", styles.navbarLogo)} {...props}>
      {logoImageUrl != null && (
        <>
          <img className="navbar__logo" src={logoImageUrl} alt={logo.alt} />
          <div className={clsx(styles.logoSubtext)}>by</div>
          <img
            className={clsx(styles.lyftLogo)}
            src={lyftLogoUrl}
            alt={logo.alt}
          />
        </>
      )}
    </div>
  );
}

function SocialMedia({ links }: { links: SocialLink[] }): JSX.Element {
  return (
    <div style={{ paddingTop: "2.5%" }}>
      {links.map((media, idx) => (
        <Link
          key={idx}
          style={{ textDecoration: "none" }}
          target="_blank"
          rel="noopener noreferrer"
          href={media.href}
        >
          {media.icon !== undefined && (
            <i className={clsx(styles.icon, media.icon)} />
          )}
        </Link>
      ))}
    </div>
  );
}

function Links({ links }: { links: FooterLink[] }): JSX.Element {
  if (links === undefined || links.length <= 0) {
    return <></>;
  }

  return (
    <div className="row footer__links">
      {links.map((linkItem, i) => (
        <div key={i} className="col footer__col">
          {linkItem.title != null ? (
            <h4 className="footer__title">{linkItem.title}</h4>
          ) : null}
          {linkItem.items != null &&
          Array.isArray(linkItem.items) &&
          linkItem.items.length > 0 ? (
            <ul className={clsx("footer__items", styles.footerAdditional)}>
              {linkItem.items.map((item, key) =>
                item.html !== undefined ? (
                  <li
                    key={key}
                    className={clsx("footer__item")}
                    dangerouslySetInnerHTML={{
                      __html: item.html,
                    }}
                  />
                ) : (
                  <li key={item.to} className="footer__item">
                    <FooterLinkItem {...item} />
                  </li>
                )
              )}
            </ul>
          ) : null}
        </div>
      ))}
    </div>
  );
}

function Footer(): JSX.Element {
  const { siteConfig } = useDocusaurusContext();
  const themeConfig = { ...siteConfig.themeConfig } as ThemeConfig;
  const { footer } = themeConfig;

  if (footer === undefined) {
    return <></>;
  }

  const { copyright, links = [] } = footer ?? {};

  const classNames = ["footer"];
  if (typeof window !== "undefined" && window.location.pathname === "/") {
    if (footer.style === "dark") {
      classNames.push(styles.gradientDark);
    } else {
      classNames.push(styles.gradient);
    }
  } else {
    classNames.push(styles.noGradient);
  }

  return (
    <footer
      className={clsx(...classNames, {
        "footer--dark": footer.style === "dark",
      })}
    >
      <div className={clsx("container", styles.container)}>
        <div className="container">
          <Logo />
          <SocialMedia links={socialLinks} />
        </div>
        <div className={clsx("container", styles.container)}>
          <Links links={links} />
        </div>
      </div>
      <div className={clsx(styles.section)}>
        {copyright !== undefined && (
          <div className={clsx("text--center", styles.copyright)}>
            <div
              dangerouslySetInnerHTML={{
                __html: copyright,
              }}
            />
            <div style={{ fontSize: ".875rem" }}>
              This site is powered by{" "}
              <a href="https://www.netlify.com/" target="blank">
                Netlify
              </a>
              .
            </div>
          </div>
        )}
      </div>
    </footer>
  );
}

export default Footer;
