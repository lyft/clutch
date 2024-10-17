import * as React from "react";
import type { LinkProps as MuiLinkProps, Theme } from "@mui/material";
import { Link as MuiLink } from "@mui/material";

import styled from "./styled";

type TextTransform = "none" | "capitalize" | "uppercase" | "lowercase" | "initial" | "inherit";

const StyledLink = styled(MuiLink)<{
  $textTransform: LinkProps["textTransform"];
  $whiteSpace: LinkProps["whiteSpace"];
}>(
  ({ theme }: { theme: Theme }) => ({
    display: "flex",
    width: "100%",
    maxWidth: "fit-content",
    fontSize: "14px",
    color: theme.palette.primary[600],
    overflow: "hidden",
    textOverflow: "ellipsis",
  }),
  props => ({
    textTransform: props.$textTransform,
    ...(props.$whiteSpace ? { whiteSpace: props.$whiteSpace } : {}),
  })
);

export interface LinkProps extends Pick<MuiLinkProps, "href" | "children"> {
  textTransform?: TextTransform;
  whiteSpace?: React.CSSProperties["whiteSpace"];
  target?: React.AnchorHTMLAttributes<HTMLAnchorElement>["target"];
}

export const Link = ({
  href,
  textTransform = "none",
  target = "_blank",
  children,
  whiteSpace,
  ...props
}: LinkProps) => (
  <StyledLink
    href={href}
    target={target}
    rel="noopener noreferrer"
    $textTransform={textTransform}
    $whiteSpace={whiteSpace}
    {...props}
  >
    {children}
  </StyledLink>
);
