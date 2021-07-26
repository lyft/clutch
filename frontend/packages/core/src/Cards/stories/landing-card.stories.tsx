import * as React from "react";
import type { Meta } from "@storybook/react";

import type { LandingCardProps } from "../landing-card";
import { LandingCard } from "../landing-card";

export default {
  title: "Core/Cards/Landing Card",
  component: LandingCard,
} as Meta;

const Template = (props: LandingCardProps) => <LandingCard {...props} />;

export const Primary = Template.bind({});
Primary.args = {
  group: "AWS",
  title: "EC2: Terminate Instance",
  description: "Search for and terminate an EC2 instance.",
};
