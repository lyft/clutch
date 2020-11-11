import React from "react";
import type { LinkProps as MuiLinkProps } from "@material-ui/core";
import { Link as MuiLink } from "@material-ui/core";
import styled from "styled-components";

const StyledLink = styled(MuiLink)`
  padding-left: 5px;
`;

export interface LinkProps extends Pick<MuiLinkProps, "href"> {}

const Link: React.FC<LinkProps> = ({ href, children }) => (
  <StyledLink href={href} target="_blank" rel="noopener noreferrer">
    {children}
  </StyledLink>
);

export default Link;
