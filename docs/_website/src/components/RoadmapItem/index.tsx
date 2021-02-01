import React from 'react';
import {MDXProvider} from '@mdx-js/react';
import MDXComponents from '@theme/MDXComponents';

import "./styles.css";

const RoadmapItem = (props) => {
    const [displayLong, setDisplayLong] = React.useState(false);

    return <div className="ri-container">
        <div className="ri-description-short">
            <div className="ri-icon"><span className="fe fe-zap" /></div>
            <div className="ri-detail">
                <div className="ri-title">{props.title}</div>
                <div className="ri-description">
                    {props.description}
                    {React.Children.count(props.children) > 0 && <span className="ri-more fe fe-more-horizontal" onClick={() => setDisplayLong(!displayLong)}/>}
                </div>
            </div>
        </div>
        {displayLong && <div className="ri-description-long"><MDXProvider components={MDXComponents}>{props.children}</MDXProvider></div>}
    </div>
}

export default RoadmapItem;