import React from "react";
import { MDXProvider } from "@mdx-js/react";
import MDXComponents from "@theme/MDXComponents";
import Link from "@docusaurus/Link";
import type { Props as LinkProps } from "@docusaurus/Link";

import "./styles.css";

interface CommunityCardProps extends Pick<LinkProps, "to"> {
  icon: string;
  children: React.ReactChildren;
}

const CommunityCard = ({
  to,
  icon,
  children,
}: CommunityCardProps): JSX.Element => (
  <Link className="cc-container" to={to}>
    <div className="cc-icon">
      <span className={`fe fe-${icon}`} />
    </div>
    <div className="cc-content">
      <MDXProvider components={MDXComponents}>{children}</MDXProvider>
    </div>
  </Link>
);

export default CommunityCard;
