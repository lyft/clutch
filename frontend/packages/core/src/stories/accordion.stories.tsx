import * as React from "react";
import type { Meta } from "@storybook/react";

import type { AccordionProps } from "../accordion";
import { Accordion, AccordionActions, AccordionDetails, AccordionDivider } from "../accordion";
import { Button } from "../button";

export default {
  title: "Core/Accordion",
  component: Accordion,
} as Meta;

const Template = (props: AccordionProps) => (
  <Accordion {...props}>
    <AccordionDetails>
      Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut
      labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco
      laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in
      voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat
      non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
    </AccordionDetails>
    <AccordionDivider />
    <AccordionActions>
      <Button text="Okay" />
    </AccordionActions>
  </Accordion>
);

export const Basic = Template.bind({});
Basic.args = {
  title: "Hello world!",
  defaultExpanded: true,
};
