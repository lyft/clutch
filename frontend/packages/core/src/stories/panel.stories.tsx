import React from "react";
import type { Meta } from "@storybook/react";

import type { ExpansionPanelProps } from "../panel";
import ExpansionPanel from "../panel";

export default {
  title: "Core/ExpansionPanel",
  component: ExpansionPanel,
} as Meta;

const Template = (props: ExpansionPanelProps) => (
  <ExpansionPanel {...props}>
    <img alt="clutch logo" src="https://clutch.sh/img/navigation/logo.svg" height="100px" />
  </ExpansionPanel>
);

export const Primary = Template.bind({});
Primary.args = {
  heading: "Check this out!",
  summary: "This is an expansion panel.",
};

export const Expanded = Template.bind({});
Expanded.args = {
  ...Primary.args,
  expanded: true,
};
