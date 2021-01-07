import React from "react";
import styled from "@emotion/styled";
import { Grid, Paper } from "@material-ui/core";
import type { Color } from "@material-ui/lab/Alert";

import { Alert } from "./alert"

const NotePanelContainer = styled(Grid)({
  "> *": {
    padding: "4px 0",
  }
});

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
    <Paper elevation={0}>
      <Alert severity={severity}>
        <Grid container justify="center" alignItems="center">
          {children}
        </Grid>
      </Alert>
    </Paper>
  );
};

const NotePanel: React.FC<NotePanelProps> = ({ direction = "column", notes, children }) => (
  <NotePanelContainer container direction={direction} justify="center" alignContent="space-between" wrap="nowrap">
    {notes?.map((note: NoteConfig) => (
      <Note key={note.text} severity={note.severity}>
        {note.text}
      </Note>
    ))}
    {children}
  </NotePanelContainer>
);

export { Note, NotePanel };
