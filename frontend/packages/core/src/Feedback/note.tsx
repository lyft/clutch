import React from "react";
import { Grid, Paper } from "@material-ui/core";
import type { Color } from "@material-ui/lab/Alert";
import MuiAlert from "@material-ui/lab/Alert";
import styled from "styled-components";

const Container = styled(Paper)`
  margin: 1%;
`;

const Alert = styled(MuiAlert)`
  align-items: center;
`;

export interface NoteProps {
  severity?: Color;
}

export interface NoteConfig extends NoteProps {
  text: string;
}

export interface NotePanelProps {
  direction?: "row" | "column";
  notes?: NoteConfig[];
}

const Note: React.FC<NoteProps> = ({ severity = "info", children }) => {
  return (
    <Container elevation={0}>
      <Alert severity={severity}>
        <Grid container justify="center" alignItems="center">
          {children}
        </Grid>
      </Alert>
    </Container>
  );
};

const NotePanel: React.FC<NotePanelProps> = ({ direction = "column", notes, children }) => (
  <Grid container direction={direction} justify="center" alignContent="space-between" wrap="nowrap">
    {notes?.map((note: NoteConfig) => (
      <Note key={note.text} severity={note.severity}>
        {note.text}
      </Note>
    ))}
    {children}
  </Grid>
);

export { Note, NotePanel };
