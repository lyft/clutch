import * as React from "react";
import type { Meta } from "@storybook/react";

import type { EmojiProps } from "../emojis";
import { SadIcon } from "../emojis";
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

export const Sad: React.FC<EmojiProps> = ({ size }) => <SadIcon size={size} />;
