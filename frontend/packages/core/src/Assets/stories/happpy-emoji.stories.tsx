import * as React from "react";
import type { Meta } from "@storybook/react";

import HappyEmoji from "../emojis/happy";
import type { SVGProps } from "../global";
import { VARIANTS } from "../global";

export default {
  title: "Core/Assets/emojis",
  component: HappyEmoji,
  argTypes: {
    size: {
      options: VARIANTS,
      control: { type: "select" },
    },
  },
} as Meta;

export const Happy: React.FC<SVGProps> = ({ size }) => <HappyEmoji size={size} />;
