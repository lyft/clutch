import React from "react";
import { Meta } from "@storybook/react/types-6-0";
import Logo from "../logo";

export default {
  title: "Core/AppLayout/Logo",
  component: Logo,
} as Meta;

export const Primary: React.FC<{}> = () => <Logo />;