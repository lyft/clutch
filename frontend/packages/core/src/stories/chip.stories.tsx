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
  label: "Neutral",
};

export const NeutralBadChip = Template.bind({});
NeutralBadChip.args = {
  variant: "NEUTRALBAD",
  label: "NeutralBad",
};

export const NeutralGoodChip = Template.bind({});
NeutralGoodChip.args = {
  variant: "NEUTRALGOOD",
  label: "NeutralGood",
};

export const RunningChip = Template.bind({});
RunningChip.args = {
  variant: "GOOD",
  label: "Good",
};

export const SucceededChip = Template.bind({});
SucceededChip.args = {
  variant: "BEST",
  label: "Best",
};
