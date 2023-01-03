import * as React from "react";
import type { LinkProps as MuiLinkProps } from "@mui/material";
import { Link as MuiLink } from "@mui/material";

import { styled } from "../Utils";

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
}

export const Link = ({ href, textTransform = "none", children, ...props }: LinkProps) => (
  <StyledLink
    href={href}
    target="_blank"
    rel="noopener noreferrer"
    $textTransform={textTransform}
    {...props}
  >
    {children}
  </StyledLink>
);
