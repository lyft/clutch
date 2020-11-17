import React from "react";
import type { LinkProps as MuiLinkProps } from "@material-ui/core";
import { Link as MuiLink } from "@material-ui/core";
import styled from "styled-components";

const StyledLink = styled(MuiLink)`
  ${({ ...props }) => `
  padding-left: 5px;
  padding-top: 10px;
  font-size: 16px;
  display: flex;
  ${props["data-max-width"] && "width: 100%;"}
  max-width: ${props["data-max-width"] || "fit-content"};
  `}
`;

export interface LinkProps extends Pick<MuiLinkProps, "href"> {
  maxWidth?: string;
}

const Link: React.FC<LinkProps> = ({ href, maxWidth, children }) => (
  <StyledLink href={href} target="_blank" rel="noopener noreferrer" data-max-width={maxWidth}>
    {children}
  </StyledLink>
);

export default Link;
