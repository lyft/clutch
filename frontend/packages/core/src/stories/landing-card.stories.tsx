import * as React from "react";
import type { Meta } from "@storybook/react";

import type { LandingCardProps } from "../card";
import { LandingCard as LandingCardComponent } from "../card";

export default {
  title: "Core/Card/Landing Card",
  component: LandingCardComponent,
} as Meta;

const Template = (props: LandingCardProps) => <LandingCardComponent {...props} />;

export const LandingCard = Template.bind({});
LandingCard.args = {
  group: "AWS",
  title: "EC2: Terminate Instance",
  description: "Search for and terminate an EC2 instance.",
};
