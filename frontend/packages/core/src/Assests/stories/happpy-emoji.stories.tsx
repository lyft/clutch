import * as React from "react";
import type { Meta } from "@storybook/react";

import type { EmojiProps } from "../emojis";
import { HappyIcon } from "../emojis";
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

export const Happy: React.FC<EmojiProps> = ({ size }) => <HappyIcon size={size} />;
