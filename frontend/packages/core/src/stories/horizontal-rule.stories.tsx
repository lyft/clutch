import * as React from "react";
import ErrorIcon from "@material-ui/icons/Error";
import type { Meta } from "@storybook/react";

import type { HorizontalRuleProps } from "../horizontal-rule";
import { HorizontalRule } from "../horizontal-rule";

export default {
  title: "Core/HorizontalRule",
  component: HorizontalRule,
} as Meta;

const Template = (props: HorizontalRuleProps) => <HorizontalRule {...props} />;

export const Basic = Template.bind({});
Basic.args = {
  children: "OR",
};

export const WithIcon = Template.bind({});
WithIcon.args = {
  children: (
    <>
      <ErrorIcon />
      &nbsp;Proceed with caution
    </>
  ),
};

export const WithoutText = Template.bind({});
WithoutText.args = {};
