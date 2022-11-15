import React, { useCallback, useState, useEffect } from "react";
import clsx from "clsx";
import Link from "@docusaurus/Link";
import useDocusaurusContext from "@docusaurus/useDocusaurusContext";
import useBaseUrl from "@docusaurus/useBaseUrl";

import ColorModeToggle from "@theme/ColorModeToggle";
import SearchBar from "@theme/SearchBar";
import { useColorMode, useWindowSize } from "@docusaurus/theme-common";
import type { NavLinkProps as BaseNavLinkProps } from "@docusaurus/Link";
import {
  useLockBodyScroll,
  // @ts-expect-error
} from "@docusaurus/theme-common/internal";

import Logo from "@theme/Logo";

import styles from "./styles.module.css";
import { ThemeConfig } from "../types";

// retrocompatible with v1
const DefaultNavItemPosition = "right";

interface ItemProps {
  to: string;
  activeBasePath: string;
  icon: string;
  label: string;
  className?: string;
}

// items defined here instead of config so they can have an associated icon
const items = [
  {
    to: "docs/about/what-is-clutch",
    activeBasePath: "docs",
    icon: "fe fe-book",
    label: "Docs",
  },
  {
    to: "blog",
    activeBasePath: "blog",
    icon: "fe fe-rss",
    label: "Blog",
  },
  {
    to: "docs/community",
    activeBasePath: "docs",
    icon: "fe fe-message-square",
    label: "Community",
  },
  {
    href: "https://github.com/lyft/clutch",
    icon: "fe fe-github",
    label: "GitHub",
  },
] as ItemProps[];

interface NavLinkProps extends BaseNavLinkProps {
  to: string;
  label: string;
  activeClassName?: string;
  className?: string;
  prependBaseUrlToHref?: boolean;
  icon: string;
}

function NavLink({
  to,
  href = "",
  label,
  activeClassName = "navbar__link--active",
  prependBaseUrlToHref = false,
  icon,
  ...props
}: NavLinkProps): JSX.Element {
  const toUrl = useBaseUrl(to);
  const normalizedHref = useBaseUrl(href, { forcePrependBaseUrl: true });

  return (
    <Link
      {...(href !== ""
        ? {
            target: "_blank",
            rel: "noopener noreferrer",
            href: prependBaseUrlToHref ? normalizedHref : href,
          }
        : {
            isNavLink: true,
            to: toUrl,
          })}
      {...props}
    >
      <span className={clsx(styles.navbarItemIcon, icon)} />
      <span className={styles.navbarItemLabel}>{label}</span>
    </Link>
  );
}

interface NavItemProps extends NavLinkProps {
  items?: ItemProps[];
  position?: "right" | "left";
  className?: string;
}

function NavItem({
  items,
  position = DefaultNavItemPosition,
  className,
  ...props
}: NavItemProps): JSX.Element {
  const navLinkClassNames = (extraClassName, isDropdownItem = false): string =>
    clsx(
      {
        "navbar__item navbar__link": !isDropdownItem,
        dropdown__link: isDropdownItem,
      },
      extraClassName
    );

  if (items == null) {
    return <NavLink className={navLinkClassNames(className)} {...props} />;
  }

  return (
    <div
      className={clsx("navbar__item", "dropdown", "dropdown--hoverable", {
        "dropdown--left": position === "left",
        "dropdown--right": position === "right",
      })}
    >
      <NavLink
        className={navLinkClassNames(className)}
        {...props}
        onClick={(e) => e.preventDefault()}
        onKeyDown={(e) => {
          if (e.key === "Enter") {
            e?.currentTarget?.parentElement?.classList.toggle("dropdown--show");
          }
        }}
      >
        {props.label}
      </NavLink>
      <ul className="dropdown__menu">
        {items.map(
          ({ className: childItemClassName, ...childItemProps }, i) => (
            <li key={i}>
              <NavLink
                activeClassName="dropdown__link--active"
                className={navLinkClassNames(childItemClassName, true)}
                {...childItemProps}
              />
            </li>
          )
        )}
      </ul>
    </div>
  );
}

interface MobileNavItemProps extends NavItemProps {
  onClick: any;
}

function MobileNavItem({
  items,
  className,
  ...props
}: MobileNavItemProps): JSX.Element {
  // Need to destructure position from props so that it doesn't get passed on.
  const navLinkClassNames = (extraClassName, isSubList = false): string =>
    clsx(
      "menu__link",
      {
        "menu__link--sublist": isSubList,
      },
      extraClassName
    );

  if (items == null) {
    return (
      <li className="menu__list-item">
        <NavLink className={navLinkClassNames(className)} {...props} />
      </li>
    );
  }

  return (
    <li className="menu__list-item">
      <NavLink className={navLinkClassNames(className, true)} {...props}>
        {props.label}
      </NavLink>
      <ul className="menu__list">
        {items.map(
          ({ className: childItemClassName, ...childItemProps }, i) => (
            <li className="menu__list-item" key={i}>
              <NavLink
                activeClassName="menu__link--active"
                className={navLinkClassNames(childItemClassName)}
                {...childItemProps}
                onClick={props.onClick}
              />
            </li>
          )
        )}
      </ul>
    </li>
  );
}

// If split links by left/right
// if position is unspecified, fallback to right (as v1)
function splitLinks(links): { leftLinks: any; rightLinks: any } {
  const leftLinks = links.filter(
    (linkItem) => (linkItem.position ?? DefaultNavItemPosition) === "left"
  );
  const rightLinks = links.filter(
    (linkItem) => (linkItem.position ?? DefaultNavItemPosition) === "right"
  );
  return {
    leftLinks,
    rightLinks,
  };
}

function Navbar(): JSX.Element {
  const { siteConfig } = useDocusaurusContext();
  const themeConfig = { ...siteConfig.themeConfig } as ThemeConfig;
  const hideOnScroll = themeConfig?.navbar?.hideOnScroll ?? false;
  const disableColorModeSwitch = themeConfig?.colorMode?.disableSwitch ?? false;
  const [sidebarShown, setSidebarShown] = useState(false);

  const { colorMode, setLightTheme, setDarkTheme } = useColorMode();

  useLockBodyScroll(sidebarShown);

  const showSidebar = useCallback(() => {
    setSidebarShown(true);
  }, [setSidebarShown]);
  const hideSidebar = useCallback(() => {
    setSidebarShown(false);
  }, [setSidebarShown]);

  const onToggleChange = useCallback(
    (e) => ((e.target.checked as boolean) ? setDarkTheme() : setLightTheme()),
    [setLightTheme, setDarkTheme]
  );

  const windowSize = useWindowSize();

  useEffect(() => {
    if (windowSize === "desktop") {
      setSidebarShown(false);
    }
  }, [windowSize]);

  const { leftLinks, rightLinks } = splitLinks(items);

  return (
    <nav
      className={clsx(
        "navbar",
        colorMode === "light" ? "navbar--light" : "navbar--dark",
        "navbar--fixed-top",
        styles.navbarCustom,
        {
          "navbar-sidebar--show": sidebarShown,
          [styles.navbarHideable]: hideOnScroll,
        }
      )}
    >
      <div className="navbar__inner">
        <div className="navbar__items">
          {items != null && items.length !== 0 && (
            <div
              aria-label="Navigation bar toggle"
              className="navbar__toggle"
              role="button"
              tabIndex={0}
              onClick={showSidebar}
              onKeyDown={showSidebar}
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="30"
                height="30"
                viewBox="0 0 30 30"
                role="img"
                focusable="false"
              >
                <title>Menu</title>
                <path
                  stroke="currentColor"
                  strokeLinecap="round"
                  strokeMiterlimit="10"
                  strokeWidth="2"
                  d="M4 7h22M4 15h22M4 23h22"
                />
              </svg>
            </div>
          )}
          <div className={clsx("navbar__brand", styles.navbarLogoCustom)}>
            <Logo
              imageClassName={clsx("navbar__logo", styles.navbarLogoCustom)}
            />
            <img
              className={clsx("navbar__title", styles.navbarLogoTextCustom)}
              src={useBaseUrl("img/navigation/logoText.svg")}
            />
          </div>
          {leftLinks.map((linkItem, i) => (
            <NavItem {...linkItem} key={i} />
          ))}
        </div>
        <div className="navbar__items navbar__items--right">
          <SearchBar />
          {rightLinks.map((linkItem, i) => (
            <NavItem {...linkItem} key={i} />
          ))}
          {!disableColorModeSwitch && (
            <ColorModeToggle
              value={colorMode}
              className={styles.displayOnlyInLargeViewport}
              onChange={onToggleChange}
            />
          )}
        </div>
      </div>
      <div
        role="presentation"
        className="navbar-sidebar__backdrop"
        onClick={hideSidebar}
      />
      <div className="navbar-sidebar">
        <div className="navbar-sidebar__brand">
          <div
            className={clsx("navbar__brand", styles.navbarLogoCustom)}
            onClick={hideSidebar}
          >
            <Logo
              imageClassName={clsx("navbar__logo", styles.navbarLogoCustom)}
            />
            <img
              className={clsx("navbar__title", styles.navbarLogoTextCustom)}
              src={useBaseUrl("img/navigation/logoText.svg")}
            />
          </div>
          {!disableColorModeSwitch && sidebarShown && (
            <ColorModeToggle value={colorMode} onChange={onToggleChange} />
          )}
        </div>
        <div className="navbar-sidebar__items">
          <div className="menu">
            <ul className="menu__list">
              {items.map((linkItem: ItemProps, i) => (
                <MobileNavItem {...linkItem} onClick={hideSidebar} key={i} />
              ))}
            </ul>
          </div>
        </div>
      </div>
    </nav>
  );
}

export default Navbar;
