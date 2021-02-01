import React from 'react';
import Link from '@docusaurus/Link';

import './styles.css';

const LinkCard = (props) => <Link className="lc-container" to={props.to}>
    <div className="lc-title">
        {props.title}
    </div>
    <div className="lc-description">
        {props.description}
    </div>
</Link>

export default LinkCard;