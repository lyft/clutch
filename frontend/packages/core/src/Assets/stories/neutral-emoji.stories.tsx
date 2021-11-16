import * as React from "react";
import type { Meta } from "@storybook/react";

import NeutralEmoji from "../emojis/neutral";
import type { SVGProps } from "../global";
import { VARIANTS } from "../global";

export default {
  title: "Core/Assets/emojis",
  component: NeutralEmoji,
  argTypes: {
    size: {
      options: VARIANTS,
      control: { type: "select" },
    },
  },
} as Meta;

export const Neutral: React.FC<SVGProps> = ({ size }) => <NeutralEmoji size={size} />;
