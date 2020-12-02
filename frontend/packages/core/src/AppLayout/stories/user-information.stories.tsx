import * as React from "react";
import type { Meta } from "@storybook/react";

import { UserInformation } from "../user";

export default {
  title: "Core/AppLayout/User Information",
  component: UserInformation,
  parameters: {
    backgrounds: {
      default: "clutch",
      values: [{ name: "clutch", value: "#131C5F" }],
    },
    layout: "centered",
  },
} as Meta;

export const Primary: React.FC<{}> = () => <UserInformation />;
