import React from "react";
import type { Meta } from "@storybook/react";

import type { NoteProps } from "../note";
import { Note, NotePanel } from "../note";

export default {
  title: "Core/Feedback/Note",
  component: Note,
} as Meta;

const Template = (props: NoteProps) => <Note {...props}>This is a note</Note>;

export const Success = Template.bind({});
Success.args = {
  severity: "success",
};

export const Error = Template.bind({});
Error.args = {
  severity: "error",
};

export const Info = Template.bind({});
Info.args = {
  severity: "info",
};

export const Warning = Template.bind({});
Warning.args = {
  severity: "warning",
};
