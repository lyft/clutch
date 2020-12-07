import React from "react";
import type { Meta } from "@storybook/react";

import type { CardProps } from "../card";
import Card from "../card";

export default {
  title: "Core/Card",
  component: Card,
} as Meta;

const Template = (props: CardProps) => <Card {...props} />;

export const Primary = Template.bind({});
Primary.args = {
  title: "Hello World",
  description: "This is a card",
};

export const CustomColor = Template.bind({});
CustomColor.argTypes = {
  backgroundColor: { control: "color" },
  titleColor: { control: "color" },
  descriptionColor: {
    control: {
      type: "select",
      options: ["textPrimary", "textSecondary"],
    },
  },
};
CustomColor.args = {
  ...Primary.args,
  backgroundColor: "#2D3F50",
  titleColor: "#ffffff",
  descriptionColor: "textSecondary",
};
