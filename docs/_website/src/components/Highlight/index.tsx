import React from 'react';


const Highlight = ({children, color}) => (
    <span
      style={{
        backgroundColor: 'var(--ifm-color-primary)',
        borderRadius: '2px',
        color: '#fff',
        padding: '0.2rem',
      }}>
      {children}
    </span>
  );

  export default Highlight;