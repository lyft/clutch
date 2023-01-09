import React from "react";
import type { Meta } from "@storybook/react";

import NotFoundComponent from "../not-found";

export default {
  title: "Core/NotFound",
  component: NotFoundComponent,
} as Meta;

const Template = () => <NotFoundComponent />;
export const NotFound = Template.bind({});
