import * as React from "react";
import type { Meta } from "@storybook/react";

import SadIcon from "../emojis/sad-emoji";
import type { SVGProps } from "../global";
import { VARIANTS } from "../global";

export default {
  title: "Core/Assets/emojis",
  component: SadIcon,
  argTypes: {
    size: {
      options: VARIANTS,
      control: { type: "select" },
    },
  },
} as Meta;

export const Sad: React.FC<SVGProps> = ({ size }) => <SadIcon size={size} />;
