import * as React from "react";

import type { SVGProps } from "../global";
import { StyledSVG } from "../global";

const HappyEmoji = ({ size, ...props }: SVGProps) => (
  <StyledSVG
    size={size}
    viewBox="0 0 48 48"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
    {...props}
  >
    <circle cx="24" cy="24" r="24" fill="#F59E0B" />
    <circle cx="21.913" cy="21.913" r="21.913" fill="#FBBF24" />
    <path
      d="M19.8261 35.4783C27.3179 35.4783 33.3913 29.8721 33.3913 22.9565H6.26086C6.26086 29.8721 12.3342 35.4783 19.8261 35.4783Z"
      fill="#0D1030"
    />
    <mask
      id="mask0_76_9148"
      style={{ maskType: "alpha" }}
      maskUnits="userSpaceOnUse"
      x="6"
      y="22"
      width="28"
      height="14"
    >
      <path
        d="M19.8261 35.4783C27.3179 35.4783 33.3913 29.8721 33.3913 22.9565H6.26086C6.26086 29.8721 12.3342 35.4783 19.8261 35.4783Z"
        fill="#0D1030"
      />
    </mask>
    <g mask="url(#mask0_76_9148)">
      <circle cx="19.8261" cy="35.4782" r="11.3043" fill="#DB3615" />
    </g>
    <ellipse
      rx="4.17391"
      ry="4.17391"
      transform="matrix(1 0 0 -1 10.4348 16.6957)"
      fill="#0D1030"
    />
    <circle r="4.17391" transform="matrix(1 0 0 -1 10.4348 18.7826)" fill="#FBBF24" />
    <ellipse
      rx="4.17391"
      ry="4.17391"
      transform="matrix(1 0 0 -1 29.2174 16.6957)"
      fill="#0D1030"
    />
    <circle r="4.17391" transform="matrix(1 0 0 -1 29.2174 18.7826)" fill="#FBBF24" />
  </StyledSVG>
);

export default HappyEmoji;
