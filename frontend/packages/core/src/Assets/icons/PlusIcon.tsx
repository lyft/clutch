import * as React from "react";

import type { SVGProps } from "../global";

import StyledSvg from "./helpers";

interface PlusIconProps {
  props?: SVGProps;
  disabled?: boolean;
}

const PlusIcon = ({ props, disabled = false }: PlusIconProps) => (
  <StyledSvg
    width="32"
    height="32"
    viewBox="0 0 32 32"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
    {...props}
  >
    <g id="Type=Primary, Style=Icon, State=Pressed">
      <path
        xmlns="http://www.w3.org/2000/svg"
        d="M0 16C0 7.16344 7.16344 0 16 0V0C24.8366 0 32 7.16344 32 16V16C32 24.8366 24.8366 32 16 32V32C7.16344 32 0 24.8366 0 16V16Z"
        fill={!disabled ? "#3548d4" : "#D3D3D3"}
      />
      <g id="Frame 1">
        <g id="Clutch Icons">
          <path id="Vector" d="M23 17H17V23H15V17H9V15H15V9H17V15H23V17Z" fill="white" />
        </g>
      </g>
    </g>
  </StyledSvg>
);

export default PlusIcon;
