import * as React from "react";
import type { Meta } from "@storybook/react";

import SadEmoji from "../emojis/sad";
import type { SVGProps } from "../global";
import { VARIANTS } from "../global";

export default {
  title: "Core/Assets/emojis",
  component: SadEmoji,
  argTypes: {
    size: {
      options: VARIANTS,
      control: { type: "select" },
    },
  },
} as Meta;

export const Sad: React.FC<SVGProps> = ({ size }) => <SadEmoji size={size} />;
