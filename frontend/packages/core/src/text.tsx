import * as React from "react";
import styled from "@emotion/styled";
import { alpha, Fab, Grid, Theme } from "@mui/material";

import { ClipboardButton } from "./button";

const CopyButtonContainer = styled(Grid)({
  marginLeft: "7px",
  flex: 0,
});

const ContentContainer = styled(Grid)({
  flex: 1,
});

const Pre = styled("pre")(({ theme }: { theme: Theme }) => ({
  border: `1px solid ${alpha(theme.palette.secondary[900], 0.38)}`,
  backgroundColor: alpha(theme.palette.secondary[900], 0.12),
  borderRadius: "4px",
  fontSize: "16px",
  color: theme.palette.primary[800],
  padding: "12px 16px",
  flex: 1,
  whiteSpace: "pre-wrap",
  wordWrap: "break-word",
  flexDirection: "row-reverse",
  display: "flex",
  overflowY: "scroll",
}));

const StyledFab = styled(Fab)(({ theme }: { theme: Theme }) => ({
  background: theme.palette.secondary[200],
  "&:hover": {
    background: theme.palette.secondary[50],
  },
}));

interface CodeProps {
  children: string;
  showCopyButton?: boolean;
}

const Code = ({ children, showCopyButton = true }: CodeProps) => (
  <Pre>
    {showCopyButton && (
      // TODO: Figure out a more permanent fix for the copy button
      <CopyButtonContainer container justifyContent="flex-end">
        <StyledFab variant="circular" size="small">
          <ClipboardButton text={children} />
        </StyledFab>
      </CopyButtonContainer>
    )}
    <ContentContainer justifyContent="flex-start" alignItems="center">
      {children}
    </ContentContainer>
  </Pre>
);

export default Code;
