import * as React from "react";
import InfoOutlinedIcon from "@material-ui/icons/InfoOutlined";
import type { Meta } from "@storybook/react";

import { Typography } from "../../typography";
import type { TooltipProps } from "../tooltip";
import { Tooltip, TooltipContainer } from "../tooltip";

export default {
  title: "Core/Feedback/Tooltip",
  component: Tooltip,
} as Meta;

const Template = (props: TooltipProps) => <Tooltip {...props} />;

export const Primary = Template.bind({});
Primary.args = {
  title: "Some helpful text",
  children: <InfoOutlinedIcon />,
  placement: "right-start",
};

const text = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Lacinia purus scelerisque imperdiet viverra erat lectus volutpat mi.
At purus nunc, lacus fermentum quam pulvinar vel, maecenas.`;

const MultilineText = () => (
  <>
    <TooltipContainer>
      <Typography variant="body4" color="#FFFFFF">
        {text}
      </Typography>
    </TooltipContainer>
    <TooltipContainer>
      <Typography variant="body4" color="#FFFFFF">
        {text}
      </Typography>
    </TooltipContainer>
  </>
);
export const MultilineTooltip = Template.bind({});
MultilineTooltip.args = {
  title: <MultilineText />,
  children: <InfoOutlinedIcon />,
  placement: "bottom",
  maxWidth: "500px",
};
