import * as React from "react";

import type { ChipProps } from "../chip";
import { Chip, CHIP_VARIANTS } from "../chip";

export default {
  title: "Core/Chip",
  component: Chip,
  argTypes: {
    variant: {
      control: {
        type: "select",
        options: CHIP_VARIANTS,
      },
    },
  },
};

const Template = (props: ChipProps) => <Chip {...props} />;

export const ErrorChip = Template.bind({});
ErrorChip.args = {
  variant: "error",
  label: "error",
};

export const WarningChip = Template.bind({});
WarningChip.args = {
  variant: "warn",
  label: "warn",
};

export const NeutralChip = Template.bind({});
NeutralChip.args = {
  variant: "neutral",
  label: "neutral",
};

export const NeutralBadChip = Template.bind({});
NeutralBadChip.args = {
  variant: "attention",
  label: "attention",
};

export const NeutralGoodChip = Template.bind({});
NeutralGoodChip.args = {
  variant: "active",
  label: "active",
};

export const RunningChip = Template.bind({});
RunningChip.args = {
  variant: "pending",
  label: "pending",
};

export const SucceededChip = Template.bind({});
SucceededChip.args = {
  variant: "success",
  label: "success",
};
