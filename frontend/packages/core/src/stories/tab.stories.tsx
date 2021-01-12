import * as React from "react";
import type { Meta } from "@storybook/react";

import type { TabProps } from "../tab";
import { Tab } from "../tab";

export default {
  title: "Core/Tab/Tab",
  component: Tab,
  argTypes: {
    onClick: { action: "onClick event" },
  },
} as Meta;

const Template = (props: TabProps) => <Tab {...props} />;

export const Primary = Template.bind({});
Primary.args = {
  label: "Tab 1",
};

export const Selected = Template.bind({});
Selected.args = {
  label: "Tab 1",
  selected: true,
};
