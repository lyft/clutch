import React from "react";
import { Button } from "../button"
import { Typography } from "../typography"
import { Divider, Paper } from "@material-ui/core";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";

export interface MarkdownPaneProps {
  /** the function that will be called when clicking the X button */
  onClose: () => void;
  /** The title that will be displayed at the top */
  title: string;
  /** The actual markdown text content that will be displayed */
  markdownText: string;
  /** whether or not to use gfm, defaults to true */
  useGithubFlavor?: boolean;
  //** width of the content, defaults to 950px */
  width?: string
}

// TODO(smonero): wrap this in a Mui Drawer
// TODO(smonero): give it a storybook with controls
// This is designed to be used within a Mui <Drawer>
const MarkdownPane = ({ onClose, title, markdownText, useGithubFlavor = true, width = "950px" }: MarkdownPaneProps) => {
  return (
    <span style={{ width: width }}>
      <Paper>
        <span
          style={{ display: "flex", justifyContent: "space-between", padding: "32px 0px 32px 0px" }}
        >
          <Typography variant="h1">{title}</Typography>
          <Button text="x" onClick={onClose} />
        </span>
        <Divider />
        <ReactMarkdown remarkPlugins={useGithubFlavor ? [remarkGfm] : []}>{markdownText}</ReactMarkdown>
      </Paper>
    </span>
  );
};

export default MarkdownPane;
