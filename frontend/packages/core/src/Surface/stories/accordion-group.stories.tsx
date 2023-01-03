import * as React from "react";
import type { Meta } from "@storybook/react";

import type { AccordionGroupProps } from "../accordion";
import { Accordion, AccordionDetails, AccordionGroup } from "../accordion";

export default {
  title: "Core/Accordion/Accordion Group",
  component: AccordionGroup,
} as Meta;

const Template = (props: AccordionGroupProps) => (
  <AccordionGroup {...props}>
    <Accordion title="First Accordion">
      <AccordionDetails>This is the first accordion.</AccordionDetails>
    </Accordion>
    <Accordion title="Second Accordion">
      <AccordionDetails>This is the second accordion.</AccordionDetails>
    </Accordion>
    <Accordion title="Third Accordion">
      <AccordionDetails>This is the third accordion.</AccordionDetails>
    </Accordion>
  </AccordionGroup>
);

export const Primary = Template.bind({});
Primary.args = {};

export const WithDefaultExpanded = Template.bind({});
WithDefaultExpanded.args = {
  defaultExpandedIdx: 0,
};
