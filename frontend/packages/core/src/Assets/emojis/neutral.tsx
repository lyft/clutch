import * as React from "react";

import type { SVGProps } from "../global";
import { StyledSVG } from "../global";

const NeutralEmoji = ({ size, ...props }: SVGProps) => (
  <StyledSVG
    size={size}
    viewBox="0 0 48 48"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
    {...props}
  >
    <circle cx="24" cy="24" r="24" fill="#F59E0B" />
    <circle cx="24" cy="21.913" r="21.913" fill="#FBBF24" />
    <circle r="4.17391" transform="matrix(1 0 0 -1 14.6087 18.7826)" fill="#0D1030" />
    <ellipse
      rx="4.17391"
      ry="4.17391"
      transform="matrix(1 0 0 -1 33.3913 18.7826)"
      fill="#0D1030"
    />
    <rect x="8.34778" y="27.1304" width="33.3913" height="2.08696" rx="1.04348" fill="#0D1030" />
  </StyledSVG>
);

export default NeutralEmoji;
