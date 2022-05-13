import * as React from "react";
import type { Meta } from "@storybook/react";

import type { SVGProps } from "../global";
import { VARIANTS } from "../global";
import ExperimentIcon from "../icons/ExperimentIcon";
import FireIcon from "../icons/FireIcon";
import GemIcon from "../icons/GemIcon";
import PlusIcon from "../icons/PlusIcon";
import RocketIcon from "../icons/RocketIcon";
import SirenIcon from "../icons/SirenIcon";
import SlackIcon from "../icons/SlackIcon";

export const AllIcons: React.FC<SVGProps> = ({ size }) => (
  <div>
    <FireIcon size={size} />
    <PlusIcon />
    <ExperimentIcon size={size} />
    <GemIcon size={size} />
    <RocketIcon size={size} />
    <SirenIcon size={size} />
    <SlackIcon size={size} />
  </div>
);
export default {
  title: "Core/Assets/Icons",
  component: AllIcons,
  argTypes: {
    size: {
      options: VARIANTS,
      control: { type: "select" },
    },
  },
} as Meta;
