import * as React from "react";
import type { Meta } from "@storybook/react";

import type { ChipProps } from "../chip";
import Chip from "../chip";

export default {
  title: "Core/Chip",
  component: Chip,
} as Meta;

const Template = (props: ChipProps) => <Chip {...props} />;

export const ErrorChip = Template.bind({});
ErrorChip.args = {
  variant: "WORST",
  label: "Error",
};

export const WarningChip = Template.bind({});
WarningChip.args = {
  variant: "BAD",
  label: "Warning",
};

export const NeutralChip = Template.bind({});
NeutralChip.args = {
  variant: "NEUTRAL",
  label: "Medium",
};

export const BlankChip = Template.bind({});
BlankChip.args = {
  variant: "BLANK",
  label: "Blank",
};

export const RunningChip = Template.bind({});
RunningChip.args = {
  variant: "GOOD",
  label: "Running",
};

export const SucceededChip = Template.bind({});
SucceededChip.args = {
  variant: "BEST",
  label: "Succeeded",
};
