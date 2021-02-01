import React from 'react';

import Link from '@docusaurus/Link';

import './styles.css';

const FrontendWorkflow = ({packageName, to, workflows}) => <Link to={to} className="fw-container">
    <div className="fw-title">{packageName}</div>
    <div className="fw-workflow-list">
        {workflows.map((v) => <li>{v}</li>)}
    </div>

</Link>

export default FrontendWorkflow;
