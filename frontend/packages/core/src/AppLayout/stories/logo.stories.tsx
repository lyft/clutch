import React from "react";
import type { Meta } from "@storybook/react";

import LogoComponent from "../logo";

export default {
  title: "Core/AppLayout/Logo",
  component: LogoComponent,
} as Meta;

export const Logo: React.FC<{}> = () => <LogoComponent />;
