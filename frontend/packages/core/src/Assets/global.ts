import styled from "@emotion/styled";

export const STYLE_MAP = {
  xsmall: {
    width: 18,
    height: 18,
  },
  small: {
    width: 24,
    height: 24,
  },
  medium: {
    width: 36,
    height: 36,
  },
  large: {
    width: 48,
    height: 48,
  },
};

export const VARIANTS = Object.keys(STYLE_MAP);

export type IconSizeVariant = "xsmall" | "small" | "medium" | "large";

export const StyledSVG = styled.svg<{ size?: IconSizeVariant }>(props => ({
  width: `${STYLE_MAP[props.size]?.width || STYLE_MAP.small.width}px`,
  height: `${STYLE_MAP[props.size]?.height || STYLE_MAP.small.height}px`,
}));

export interface SVGProps {
  size?: IconSizeVariant;
}
