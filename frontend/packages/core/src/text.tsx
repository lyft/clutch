import * as React from "react";
import styled from "@emotion/styled";
import { Fab, Grid } from "@material-ui/core";

import { ClipboardButton } from "./button";

const CopyButtonContainer = styled(Grid)({
  marginLeft: "7px",
  flex: 0,
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
      <CopyButtonContainer container justify="flex-end">
        <Fab variant="round" size="small">
          <ClipboardButton text={children} />
        </Fab>
      </CopyButtonContainer>
    )}
    <Grid style={{ flex: 1 }}container justify="flex-start" alignItems="center">
      {children}
    </Grid>
  </Pre>
);

export default Code;
