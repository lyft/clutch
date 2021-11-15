import * as React from "react";
import type { Meta } from "@storybook/react";

import type { EmojiProps } from "../emojis";
import { NeutralIcon } from "../emojis";
import { VARIANTS } from "../global";

export default {
  title: "Core/Assets/emojis",
  component: NeutralIcon,
  argTypes: {
    size: {
      options: VARIANTS,
      control: { type: "select" },
    },
  },
} as Meta;

export const Neutral: React.FC<EmojiProps> = ({ size }) => <NeutralIcon size={size} />;
