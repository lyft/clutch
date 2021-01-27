import React, { useCallback, useState, useEffect } from 'react';
import clsx from 'clsx';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import useBaseUrl from '@docusaurus/useBaseUrl';

import SearchBar from '@theme/SearchBar';
import Toggle from '@theme/Toggle';
import useThemeContext from '@theme/hooks/useThemeContext';
import useHideableNavbar from '@theme/hooks/useHideableNavbar';
import useLockBodyScroll from '@theme/hooks/useLockBodyScroll';
import useWindowSize, { windowSizes } from '@theme/hooks/useWindowSize';
import useLogo from '@theme/hooks/useLogo';

import styles from './styles.module.css';

// retrocompatible with v1
const DefaultNavItemPosition = 'right';

// items defined here instead of config so they can have an associated icon
var items = [
  {
    to: 'docs/about/what-is-clutch',
    activeBasePath: 'docs',
    icon: "fe fe-book",
    label: 'Docs',
  },
  {
    to: 'blog',
    activeBasePath: 'blog',
    icon: "fe fe-rss",
    label: 'Blog',
  },
  {
    to: 'docs/community',
    activeBasePath: 'docs',
    icon: "fe fe-message-square",
    label: 'Community',
  },
  {
    to: 'https://storybook.clutch.sh',
    icon: "fe fe-layers",
    label: 'Storybook',
  },
  {
    href: 'https://github.com/lyft/clutch',
    icon: "fe fe-github",
    label: 'GitHub',
  },
];

function NavLink({
  activeBasePath,
  activeBaseRegex,
  to,
  href,
  label,
  activeClassName = 'navbar__link--active',
  prependBaseUrlToHref,
  icon,
  ...props
}) {
  const toUrl = useBaseUrl(to);
  const normalizedHref = useBaseUrl(href, { forcePrependBaseUrl: true });

  return (
    <Link
      {...(href
        ? {
          target: '_blank',
          rel: 'noopener noreferrer',
          href: prependBaseUrlToHref ? normalizedHref : href,
        }
        : {
          isNavLink: true,
          to: toUrl,
        })}
      {...props}>
      <span className={clsx(styles.navbarItemIcon, icon)} />
      <span className={styles.navbarItemLabel}>{label}</span>
    </Link>
  );
}

function NavItem({
  items,
  position = DefaultNavItemPosition,
  className,
  ...props
}) {
  const navLinkClassNames = (extraClassName, isDropdownItem = false) =>
    clsx(
      {
        'navbar__item navbar__link': !isDropdownItem,
        dropdown__link: isDropdownItem,
      },
      extraClassName,
    );

  if (!items) {
    return <NavLink className={navLinkClassNames(className)} {...props} />;
  }

  return (
    <div
      className={clsx('navbar__item', 'dropdown', 'dropdown--hoverable', {
        'dropdown--left': position === 'left',
        'dropdown--right': position === 'right',
      })}>
      <NavLink
        className={navLinkClassNames(className)}
        {...props}
        onClick={(e) => e.preventDefault()}
        onKeyDown={(e) => {
          if (e.key === 'Enter') {
            e.target.parentNode.classList.toggle('dropdown--show');
          }
        }}>
        {props.label}
      </NavLink>
      <ul className="dropdown__menu">
        {items.map(({ className: childItemClassName, ...childItemProps }, i) => (
          <li key={i}>
            <NavLink
              activeClassName="dropdown__link--active"
              className={navLinkClassNames(childItemClassName, true)}
              {...childItemProps}
            />
          </li>
        ))}
      </ul>
    </div>
  );
}

function MobileNavItem({ items, position: _position, className, ...props }) {
  // Need to destructure position from props so that it doesn't get passed on.
  const navLinkClassNames = (extraClassName, isSubList = false) =>
    clsx(
      'menu__link',
      {
        'menu__link--sublist': isSubList,
      },
      extraClassName,
    );

  if (!items) {
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
        {items.map(({ className: childItemClassName, ...childItemProps }, i) => (
          <li className="menu__list-item" key={i}>
            <NavLink
              activeClassName="menu__link--active"
              className={navLinkClassNames(childItemClassName)}
              {...childItemProps}
              onClick={props.onClick}
            />
          </li>
        ))}
      </ul>
    </li>
  );
}

// If split links by left/right
// if position is unspecified, fallback to right (as v1)
function splitLinks(links) {
  const leftLinks = links.filter(
    (linkItem) => (linkItem.position ?? DefaultNavItemPosition) === 'left',
  );
  const rightLinks = links.filter(
    (linkItem) => (linkItem.position ?? DefaultNavItemPosition) === 'right',
  );
  return {
    leftLinks,
    rightLinks,
  };
}

function Navbar() {
  const {
    siteConfig: {
      themeConfig: {
        navbar: { title, hideOnScroll = false } = {},
        colorMode: {disableSwitch: disableColorModeSwitch = false} = {},
      },
    },
    isClient,
  } = useDocusaurusContext();
  const [sidebarShown, setSidebarShown] = useState(false);
  const [isSearchBarExpanded, setIsSearchBarExpanded] = useState(false);

  const { isDarkTheme, setLightTheme, setDarkTheme } = useThemeContext();
  const { navbarRef, isNavbarVisible } = useHideableNavbar(hideOnScroll);
  const { logoLink, logoLinkProps, logoImageUrl, logoAlt } = useLogo();

  useLockBodyScroll(sidebarShown);

  const showSidebar = useCallback(() => {
    setSidebarShown(true);
  }, [setSidebarShown]);
  const hideSidebar = useCallback(() => {
    setSidebarShown(false);
  }, [setSidebarShown]);

  const onToggleChange = useCallback(
    (e) => (e.target.checked ? setDarkTheme() : setLightTheme()),
    [setLightTheme, setDarkTheme],
  );

  const windowSize = useWindowSize();

  useEffect(() => {
    if (windowSize === windowSizes.desktop) {
      setSidebarShown(false);
    }
  }, [windowSize]);

  const { leftLinks, rightLinks } = splitLinks(items);

  return (
    <nav
      ref={navbarRef}
      className={clsx('navbar', 'navbar--light', 'navbar--fixed-top', styles.navbarCustom, {
        'navbar-sidebar--show': sidebarShown,
        [styles.navbarHideable]: hideOnScroll,
        [styles.navbarHidden]: !isNavbarVisible,
      })}>
      <div className="navbar__inner">
        <div className="navbar__items">
          {items != null && items.length !== 0 && (
            <div
              aria-label="Navigation bar toggle"
              className="navbar__toggle"
              role="button"
              tabIndex={0}
              onClick={showSidebar}
              onKeyDown={showSidebar}>
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="30"
                height="30"
                viewBox="0 0 30 30"
                role="img"
                focusable="false">
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
          <Link className={clsx("navbar__brand", styles.navbarLogoCustom)} to={logoLink} {...logoLinkProps}>
            {logoImageUrl != null && (
              <img
                key={isClient}
                className={clsx("navbar__logo")}
                src={logoImageUrl}
                alt={logoAlt}
              />
            )}
            <img className={clsx('navbar__title', styles.navbarLogoTextCustom, {[styles.hideLogoText]: isSearchBarExpanded})} src={useBaseUrl("img/navigation/logoText.svg")} />
          </Link>
          {leftLinks.map((linkItem, i) => (
            <NavItem {...linkItem} key={i} />
          ))}
        </div>
        <div className="navbar__items navbar__items--right">
          <SearchBar
            handleSearchBarToggle={setIsSearchBarExpanded}
            isSearchBarExpanded={isSearchBarExpanded}
          />
          {rightLinks.map((linkItem, i) => (
            <NavItem {...linkItem} key={i} />
          ))}
          {!disableColorModeSwitch && (
            <Toggle
              className={styles.displayOnlyInLargeViewport}
              aria-label="Dark mode toggle"
              checked={isDarkTheme}
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
          <Link
            className="navbar__brand"
            onClick={hideSidebar}
            to={logoLink}
            {...logoLinkProps}>
            {logoImageUrl != null && (
              <img
                key={isClient}
                className="navbar__logo"
                src={logoImageUrl}
                alt={logoAlt}
              />
            )}
            <img className={clsx('navbar__title', styles.navbarLogoTextCustom, {[styles.hideLogoText]: isSearchBarExpanded})} src={useBaseUrl("img/navigation/logoText.svg")} />
          </Link>
          {!disableColorModeSwitch && sidebarShown && (
            <Toggle
              aria-label="Dark mode toggle in sidebar"
              checked={isDarkTheme}
              onChange={onToggleChange}
            />
          )}
        </div>
        <div className="navbar-sidebar__items">
          <div className="menu">
            <ul className="menu__list">
              {items.map((linkItem, i) => (
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
