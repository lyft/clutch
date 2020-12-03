import * as React from "react";
import type { Meta } from "@storybook/react";

import Notifications from "../notifications";

export default {
  title: "Core/AppLayout/Notifications",
  component: Notifications,
  parameters: {
    backgrounds: {
      default: "header blue",
      values: [{ name: "header blue", value: "#131C5F" }],
    },
  },
} as Meta;

export const Primary: React.FC<{}> = () => <Notifications />;
