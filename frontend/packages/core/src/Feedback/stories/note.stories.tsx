import React from "react";
import type { Meta } from "@storybook/react";

import type { NoteProps } from "../note";
import { Note as NoteComponent } from "../note";

const SEVERITIES = ["error", "info", "success", "warning"];

export default {
  title: "Core/Feedback/Single",
  component: NoteComponent,
  argTypes: {
    severity: {
      options: SEVERITIES,
      control: {
        type: "select",
      },
    },
  },
} as Meta;

const Template = (props: NoteProps) => <NoteComponent {...props}>This is a note</NoteComponent>;

export const Single = Template.bind({});
Single.args = {
  severity: "success",
};
