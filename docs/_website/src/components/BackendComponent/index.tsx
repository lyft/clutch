import React from "react";

import Link from "@docusaurus/Link";
import type { Props as LinkProps } from "@docusaurus/Link";

import "./styles.css";

interface BackendComponentProps extends Pick<LinkProps, "to"> {
  name: string;
  desc: string;
}

const BackendComponent = ({
  name,
  desc,
  to,
}: BackendComponentProps): JSX.Element => (
  <Link className="bc-container" to={to} target="_blank">
    <div className="bc-title">{name}</div>
    {desc !== "" && <div className="bc-description">{desc}</div>}
  </Link>
);

export default BackendComponent;
