import React from "react";
import type { Meta } from "@storybook/react";

import { ClipboardButton } from "../../button";
import Link from "../../link";
import type { NotePanelProps } from "../note";
import { Note, NotePanel } from "../note";

export default {
  title: "Core/Feedback/Note/Panel",
  component: NotePanel,
} as Meta;

const Template = (props: NotePanelProps) => <NotePanel {...props} />;

export const Primary = Template.bind({});
Primary.args = {
  notes: [
    {
      severity: "info",
      text: "Information message.",
    },
    {
      severity: "warning",
      text: "Warning message.",
    },
  ],
};

const CompositeTemplate = (props: NotePanelProps) => {
  const docsUrl = "https://clutch.sh/docs";
  return (
    <NotePanel {...props}>
      <Note severity="info">
        Additional information can be found in the
        <Link href="https://clutch.sh/docs">Clutch documentation</Link>
      </Note>
      <Note severity="warning">
        If the link above does not work you can copy the URL: {docsUrl}
        <ClipboardButton text={docsUrl} size="small" />
      </Note>
    </NotePanel>
  );
};

export const WithChildren = CompositeTemplate.bind({});
