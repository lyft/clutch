import React from "react";

import Link from "@docusaurus/Link";
import type { Props as LinkProps } from "@docusaurus/Link";

import "./styles.css";

interface FrontendWorkflowProps extends Pick<LinkProps, "to"> {
  packageName: string;
  workflows?: React.ReactNode[];
}

const FrontendWorkflow = ({
  packageName,
  to,
  workflows,
}: FrontendWorkflowProps): JSX.Element => (
  <Link to={to} className="fw-container">
    <div className="fw-title">{packageName}</div>
    {workflows !== null && (
      <div className="fw-workflow-list">
        {workflows?.map((v, idx) => (
          <li key={idx}>{v}</li>
        ))}
      </div>
    )}
  </Link>
);

export default FrontendWorkflow;
