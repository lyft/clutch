import * as React from "react";
import type { LinkProps as MuiLinkProps } from "@mui/material";
import { Link as MuiLink } from "@mui/material";

import styled from "./styled";

type TextTransform = "none" | "capitalize" | "uppercase" | "lowercase" | "initial" | "inherit";

const StyledLink = styled(MuiLink)<{
  $textTransform: LinkProps["textTransform"];
}>(
  {
    display: "flex",
    width: "100%",
    maxWidth: "fit-content",
    fontSize: "14px",
    color: "#3548D4",
    overflow: "hidden",
    textOverflow: "ellipsis",
  },
  props => ({
    textTransform: props.$textTransform,
  })
);

export interface LinkProps extends Pick<MuiLinkProps, "href" | "children"> {
  textTransform?: TextTransform;
  target?: React.AnchorHTMLAttributes<HTMLAnchorElement>["target"];
}

export const Link = ({
  href,
  textTransform = "none",
  target = "_blank",
  children,
  ...props
}: LinkProps) => (
  <StyledLink
    href={href}
    target={target}
    rel="noopener noreferrer"
    $textTransform={textTransform}
    {...props}
  >
    {children}
  </StyledLink>
);
