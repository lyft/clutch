import React from "react";

import "./styles.css";

interface ImageProps {
  variant: string;
}

const Image = ({ variant, ...props }: ImageProps): JSX.Element => (
  <div className={variant}>
    <img {...props} />
  </div>
);

export default Image;
