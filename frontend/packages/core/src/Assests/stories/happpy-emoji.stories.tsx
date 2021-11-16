import * as React from "react";
import type { Meta } from "@storybook/react";

import HappyIcon from "../emojis/happy-emoji";
import type { SVGProps } from "../global";
import { VARIANTS } from "../global";

export default {
  title: "Core/Assets/emojis",
  component: HappyIcon,
  argTypes: {
    size: {
      options: VARIANTS,
      control: { type: "select" },
    },
  },
} as Meta;

export const Happy: React.FC<SVGProps> = ({ size }) => <HappyIcon size={size} />;
