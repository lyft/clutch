import React from "react";
import { Grid, Paper } from "@mui/material";
import type { AlertColor as Color } from "@mui/material/Alert";

import { Link } from "../link";
import { styled } from "../Utils";

import { Alert } from "./alert";

const NotePanelContainer = styled(Grid)({
  "> *": {
    padding: "4px 0",
  },
});

export interface NoteProps {
  severity?: Color;
  // Use this if you want to add a link after your text blob (useful for
  // linking to docs or a resource described in the text blob)
  link?: string;
}

export interface NoteConfig extends NoteProps {
  text: string;
  // Use this to specify the location in regards to a wizard for the note.
  // For example, you can specify it to be in the "intro" step, then in the
  // intro you iterate through your notes, selectively showing the ones that are
  // specified for that step. This allows users to have different notes
  // for different stages of a workflow.
  location?: string;
}

export interface NotePanelProps {
  direction?: "row" | "column";
  notes?: NoteConfig[];
}

const Note: React.FC<NoteProps> = ({ severity = "info", link = "", children }) => {
  return (
    <Paper elevation={0}>
      <Alert severity={severity}>
        <Grid container justifyContent="flex-start" alignItems="center">
          {children}
        </Grid>
        {link && <Link href={link}>{link}</Link>}
      </Alert>
    </Paper>
  );
};

const NotePanel: React.FC<NotePanelProps> = ({ direction = "column", notes, children }) => (
  <NotePanelContainer
    container
    direction={direction}
    justifyContent="center"
    alignContent="space-between"
    wrap="nowrap"
  >
    {notes?.map((note: NoteConfig) => (
      <Note key={note.text} severity={note.severity} link={note.link}>
        {note.text}
      </Note>
    ))}
    {children}
  </NotePanelContainer>
);

export { Note, NotePanel };
