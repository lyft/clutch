import React from "react";
import { Grid, Paper } from "@material-ui/core";
import type { Color } from "@material-ui/lab/Alert";
import Alert from "@material-ui/lab/Alert";
import styled from "styled-components";

const Container = styled(Paper)`
  ${({ ...props }) => `
  margin: 1%;
  width: 100%;
  max-width: ${props["data-max-width"]};
  `}
`;

interface NoteProps {
  severity?: Color;
  maxWidth?: string;
  direction?: "row" | "column";
}

export interface NoteConfig extends NoteProps {
  text: string;
}

interface NotePanelProps {
  direction?: "row" | "column";
  notes: NoteConfig[];
}

const NotePanel: React.FC<NotePanelProps> = ({ direction = "column", notes }) => {
  const isColumnLayout = direction === "column";
  const maxWidth = isColumnLayout ? "100%" : "300px";
  const noteDirection = isColumnLayout ? "row" : "column";
  return (
    <Grid
      container
      direction={direction}
      justify="center"
      alignContent="space-between"
      wrap="nowrap"
    >
      {notes.map((note: NoteConfig) => (
        <Note
          key={note.text}
          severity={note.severity}
          maxWidth={maxWidth}
          direction={noteDirection}
        >
          {note.text}
        </Note>
      ))}
    </Grid>
  );
};

const Note: React.FC<NoteProps> = ({ severity = "info", maxWidth, direction, children }) => {
  return (
    <Container elevation={0} data-max-width={maxWidth}>
      <Grid container justify="center" alignItems="center" direction={direction}>
        <Alert severity={severity}>{children}</Alert>
      </Grid>
    </Container>
  );
};

export { Note, NotePanel };
