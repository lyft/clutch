import React from "react";
import type { Meta } from "@storybook/react";

import Logo from "../logo";

export default {
  title: "Core/AppLayout/Logo",
  component: Logo,
} as Meta;

export const Primary: React.FC<{}> = () => <Logo />;
