import React from "react";
import type { LinkProps as MuiLinkProps } from "@material-ui/core";
import { Link as MuiLink } from "@material-ui/core";
import styled from "styled-components";

const StyledLink = styled(MuiLink)`
  ${({ ...props }) => `
  display: flex;
  ${props["data-max-width"] && "width: 100%;"}
  max-width: ${props["data-max-width"] || "fit-content"};
  font-size: ${props["data-font-size"] || ""};
  padding: ${props["data-padding"] || "0 0 0 5px"};
  `}
`;

export interface LinkProps extends Pick<MuiLinkProps, "href"> {
  maxWidth?: string;
  fontSize?: string;
  padding?: string;
}

const Link: React.FC<LinkProps> = ({ href, maxWidth, fontSize, padding, children }) => (
  <StyledLink
    href={href}
    target="_blank"
    rel="noopener noreferrer"
    data-max-width={maxWidth}
    data-font-size={fontSize}
    data-padding={padding}
  >
    {children}
  </StyledLink>
);

export default Link;
