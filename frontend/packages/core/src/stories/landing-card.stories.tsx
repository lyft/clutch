import * as React from "react";
import type { Meta } from "@storybook/react";

import type { LandingCardProps } from "../card";
import { LandingCard } from "../card";

export default {
  title: "Core/Card/Landing Card",
  component: LandingCard,
} as Meta;

const Template = (props: LandingCardProps) => <LandingCard {...props} />;

export const Primary = Template.bind({});
Primary.args = {
  group: "AWS",
  title: "EC2: Terminate Instance",
  description: "Search for and terminate an EC2 instance.",
};
