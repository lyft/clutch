import { styled } from "../Utils";

const XSMALL = 18;
const SMALL = 24;
const MEDIUM = 36;
const LARGE = 48;

export const STYLE_MAP = {
  xsmall: {
    size: XSMALL,
  },
  small: {
    size: SMALL,
  },
  medium: {
    size: MEDIUM,
  },
  large: {
    size: LARGE,
  },
};

export const VARIANTS = Object.keys(STYLE_MAP);

export type IconSizeVariant = keyof typeof STYLE_MAP;

export const StyledSVG = styled("svg")<{ size?: IconSizeVariant }>(props => ({
  width: `${STYLE_MAP[props.size]?.size || STYLE_MAP.small.size}px`,
  height: `${STYLE_MAP[props.size]?.size || STYLE_MAP.small.size}px`,
}));

export interface SVGProps extends React.SVGProps<SVGSVGElement> {
  size?: IconSizeVariant;
}
