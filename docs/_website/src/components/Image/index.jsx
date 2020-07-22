import React from 'react';

import './styles.css';

const Image = ({variant, ...props}) => <div className={variant}>
    <img {...props} />
</div>

export default Image;