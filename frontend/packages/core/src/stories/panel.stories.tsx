import React from "react";
import type { Meta } from "@storybook/react";

import type { AccordionProps } from "../panel";
import AccordionComponent from "../panel";

export default {
  title: "Core/Accordion",
  component: AccordionComponent,
} as Meta;

const Template = (props: AccordionProps) => (
  <AccordionComponent {...props}>
    <img alt="clutch logo" src="https://clutch.sh/img/navigation/logo.svg" height="100px" />
  </AccordionComponent>
);

export const Primary = Template.bind({});
Primary.args = {
  heading: "Check this out!",
  summary: "This is an expansion panel.",
};
