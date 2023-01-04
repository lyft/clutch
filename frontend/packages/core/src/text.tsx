import * as React from "react";
import styled from "@emotion/styled";
import { Fab, Grid } from "@mui/material";

import { ClipboardButton } from "./button";

const CopyButtonContainer = styled(Grid)({
  marginLeft: "7px",
  flex: 0,
});

const ContentContainer = styled(Grid)({
  flex: 1,
});

const Pre = styled.pre({
  border: "1px solid rgba(13, 16, 48, 0.38)",
  backgroundColor: "rgba(13,16,48,0.12)",
  borderRadius: "4px",
  fontSize: "16px",
  color: "#242b8c",
  padding: "12px 16px",
  flex: 1,
  whiteSpace: "pre-wrap",
  wordWrap: "break-word",
  flexDirection: "row-reverse",
  display: "flex",
  overflowY: "scroll",
});

interface CodeProps {
  children: string;
  showCopyButton?: boolean;
}

const Code = ({ children, showCopyButton = true }: CodeProps) => (
  <Pre>
    {showCopyButton && (
      // TODO: Figure out a more permanent fix for the copy button
      <CopyButtonContainer container justifyContent="flex-end">
        <Fab variant="circular" size="small">
          <ClipboardButton text={children} />
        </Fab>
      </CopyButtonContainer>
    )}
    <ContentContainer justifyContent="flex-start" alignItems="center">
      {children}
    </ContentContainer>
  </Pre>
);

export default Code;
