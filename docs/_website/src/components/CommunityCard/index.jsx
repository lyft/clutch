import React from 'react';
import {MDXProvider} from '@mdx-js/react';
import MDXComponents from '@theme/MDXComponents';
import Link from '@docusaurus/Link';


import './styles.css';

const CommunityCard = (props) => <Link className="cc-container" to={props.to}>
    <div className="cc-icon">
        <span className={`fe fe-${ props.icon }`} />
    </div>
    <div className="cc-content">
    <MDXProvider components={MDXComponents}>{props.children}</MDXProvider>
    </div>
</Link>

export default CommunityCard;