import * as React from "react";
import styled from "@emotion/styled";
import type { LinkProps as MuiLinkProps } from "@material-ui/core";
import { Link as MuiLink } from "@material-ui/core";

export const StyledLink = styled(MuiLink)(
  {
    display: "flex",
    width: "100%",
    maxWidth: "fit-content",
    fontSize: "14px",
    color: "#3548D4",
  },
  props => ({
    textTransform: props["data-text-transform"],
  })
);

export interface LinkProps extends Pick<MuiLinkProps, "href" | "children"> {
  textTransform?: "none" | "capitalize" | "uppercase" | "lowercase" | "initial" | "inherit";
}

export const Link = ({ href, textTransform = "none", children }: LinkProps) => (
  <StyledLink
    href={href}
    target="_blank"
    rel="noopener noreferrer"
    data-text-transform={textTransform}
  >
    {children}
  </StyledLink>
);
