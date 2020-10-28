import React from "react";
import { Grid } from "@material-ui/core";
import type { Meta } from "@storybook/react";

import type { LoadableProps } from "../loading";
import Loadable from "../loading";

export default {
  title: "Core/Loading",
  component: Loadable,
} as Meta;

const Template = (props: LoadableProps) => (
  <Loadable {...props}>
    <div>Hello World!</div>
  </Loadable>
);

const LargeTemplate = (props: LoadableProps) => (
  <Loadable {...props}>
    <Grid container alignItems="center">
      <img alt="clutch logo" src="https://clutch.sh/img/navigation/logo.svg" />
    </Grid>
  </Loadable>
);

export const Primary = Template.bind({});
Primary.args = {
  isLoading: true,
};

export const Overlay = LargeTemplate.bind({});
Overlay.args = {
  ...Primary.args,
  overlay: true,
};
