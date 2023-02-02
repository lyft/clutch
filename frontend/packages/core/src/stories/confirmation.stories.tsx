import React from "react";
import type { Meta } from "@storybook/react";

import ConfirmationComponent from "../confirmation";

export default {
  title: "Core/Confirmation",
  component: ConfirmationComponent,
} as Meta;

const Template = (props: { action: string }) => <ConfirmationComponent {...props} />;

export const Confirmation = Template.bind({});
Confirmation.args = {
  action: "Deletion",
};
