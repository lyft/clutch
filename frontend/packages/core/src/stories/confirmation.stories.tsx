import React from "react";
import type { Meta } from "@storybook/react";

import Confirmation from "../confirmation";

export default {
  title: "Core/Confirmation",
  component: Confirmation,
} as Meta;

const Template = (props: { action: string }) => <Confirmation {...props} />;

export const Default = Template.bind({});
Default.args = {
  action: "Deletion",
};
