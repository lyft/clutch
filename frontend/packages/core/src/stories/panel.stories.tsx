import React from "react";
import type { Meta } from "@storybook/react";

import type { AccordionProps } from "../panel";
import Accordion from "../panel";

export default {
  title: "Core/Accordion",
  component: Accordion,
} as Meta;

const Template = (props: AccordionProps) => (
  <Accordion {...props}>
    <img alt="clutch logo" src="https://clutch.sh/img/navigation/logo.svg" height="100px" />
  </Accordion>
);

export const Primary = Template.bind({});
Primary.args = {
  heading: "Check this out!",
  summary: "This is an expansion panel.",
};

export const Expanded = Template.bind({});
Expanded.args = {
  ...Primary.args,
  expanded: true,
};
