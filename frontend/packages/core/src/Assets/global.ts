import styled from "@emotion/styled";

const XSMALL = 18;
const SMALL = 24;
const MEDIUM = 36;
const LARGE = 48;

export const STYLE_MAP = {
  xsmall: {
    width: XSMALL,
    height: XSMALL,
  },
  small: {
    width: SMALL,
    height: SMALL,
  },
  medium: {
    width: MEDIUM,
    height: MEDIUM,
  },
  large: {
    width: LARGE,
    height: LARGE,
  },
};

export const VARIANTS = Object.keys(STYLE_MAP);

export type IconSizeVariant = "xsmall" | "small" | "medium" | "large";

export const StyledSVG = styled.svg<{ size?: IconSizeVariant }>(props => ({
  width: `${STYLE_MAP[props.size]?.width || STYLE_MAP.small.width}px`,
  height: `${STYLE_MAP[props.size]?.height || STYLE_MAP.small.height}px`,
}));

export interface SVGProps extends React.SVGProps<SVGSVGElement> {
  size?: IconSizeVariant;
}
