import * as React from "react";
import styled from "@emotion/styled";

const REGULAR = 400;
const MEDIUM = 500;
const BOLD = 700;

const STYLE_MAP = {
  h1: {
    size: "36",
    weight: MEDIUM,
    lineHeight: "44",
  },
  h2: {
    size: "26",
    weight: BOLD,
    lineHeight: "32",
  },
  h3: {
    size: "22",
    weight: BOLD,
    lineHeight: "28",
  },
  h4: {
    size: "20",
    weight: BOLD,
    lineHeight: "24",
  },
  h5: {
    size: "16",
    weight: BOLD,
    lineHeight: "20",
  },
  h6: {
    size: "14",
    weight: BOLD,
    lineHeight: "18",
  },
  subtitle1: {
    size: "20",
    weight: MEDIUM,
    lineHeight: "24",
  },
  subtitle2: {
    size: "16",
    weight: MEDIUM,
    lineHeight: "20",
  },
  subtitle3: {
    size: "14",
    weight: MEDIUM,
    lineHeight: "18",
  },
  body1: {
    size: "20",
    weight: REGULAR,
    lineHeight: "26",
  },
  body2: {
    size: "16",
    weight: REGULAR,
    lineHeight: "22",
  },
  body3: {
    size: "12",
    weight: REGULAR,
    lineHeight: "20",
  },
  caption1: {
    size: "16",
    weight: BOLD,
    lineHeight: "20",
    props: {
      textTransform: "uppercase",
    },
  },
  caption2: {
    size: "12",
    weight: BOLD,
    lineHeight: "16",
    props: {
      textTransform: "uppercase",
    },
  },
  overline: {
    size: "10",
    weight: REGULAR,
    lineHeight: "10",
    props: {
      textTransform: "uppercase",
      letterSpacing: "1.5px",
    },
  },
  input: {
    size: "14",
    weight: REGULAR,
    lineHeight: "18",
  },
};

export const VARIANTS = Object.keys(STYLE_MAP);

type TextVariant =
  | "h1"
  | "h2"
  | "h3"
  | "h4"
  | "h5"
  | "h6"
  | "subtitle1"
  | "subtitle2"
  | "subtitle3"
  | "body1"
  | "body2"
  | "body3"
  | "caption1"
  | "caption2"
  | "overline"
  | "input";

const StyledTypography = styled.div<{ variant: TextVariant }>(
  {
    color: "#0D1030",
  },
  props => ({
    fontSize: `${STYLE_MAP[props?.variant || 0].size}px`,
    fontWeight: STYLE_MAP[props?.variant || 0].weight,
    lineHeight: `${STYLE_MAP[props?.variant || 0].lineHeight}px`,
    ...STYLE_MAP[props?.variant || 0].props,
  })
);

export interface TypographyProps {
  variant: TextVariant;
  children: React.ReactNode;
}

const Typography = ({ variant, children }: TypographyProps) => {
  return <StyledTypography variant={variant}>{children}</StyledTypography>;
};

export { StyledTypography, Typography };
