import React from 'react';

import Link from '@docusaurus/Link';

import './styles.css';

const BackendComponent = ({name, desc, to}) => <Link className="bc-container" to={to} target="_blank">
    <div className="bc-title">{name}</div>
    {desc && <div className="bc-description">{desc}</div>}
</Link>

export default BackendComponent;