import * as React from "react";
import type { Meta } from "@storybook/react";

import NeutralIcon from "../emojis/neutral-emoji";
import type { SVGProps } from "../global";
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

export const Neutral: React.FC<SVGProps> = ({ size }) => <NeutralIcon size={size} />;
