import * as React from "react";
import type { Meta } from "@storybook/react";

import type { AccordionGroupProps } from "../accordion";
import {
  Accordion,
  AccordionDetails,
  AccordionGroup as AccordionGroupComponent,
} from "../accordion";

export default {
  title: "Core/Accordion/Accordion Group",
  component: AccordionGroupComponent,
  argTypes: {
    defaultExpandedIdx: {
      description: "The index of the tab expanded by default on mount.",
    },
  },
} as Meta;

const Template = (props: AccordionGroupProps) => (
  <AccordionGroupComponent {...props}>
    <Accordion title="First Accordion">
      <AccordionDetails>This is the first accordion.</AccordionDetails>
    </Accordion>
    <Accordion title="Second Accordion">
      <AccordionDetails>This is the second accordion.</AccordionDetails>
    </Accordion>
    <Accordion title="Third Accordion">
      <AccordionDetails>This is the third accordion.</AccordionDetails>
    </Accordion>
  </AccordionGroupComponent>
);

export const AccordionGroup = Template.bind({});
AccordionGroup.args = {};
