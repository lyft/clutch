// function Logo({...props}) {
//   const {logoLink, logoImageUrl, logoAlt} = useLogo();
//   const lyftLogoUrl = useBaseUrl('img/microsite/lyft-logo.svg');

//   return (
//     <Link className={classnames("navbar__brand", styles.navbarLogo)} to={logoLink} {...props}>
//       {logoImageUrl != null && (
//         <>
//           <img
//             className="navbar__logo"
//             src={logoImageUrl}
//             alt={logoAlt}
//           />
//           <div className={classnames(styles.logoSubtext)}>
//             by
//           </div>
//           <img
//             className={classnames(styles.lyftLogo)}
//             src={lyftLogoUrl}
//             alt={logoAlt}
//           />
//         </>
//       )}

//     </Link>
//   );
// };
