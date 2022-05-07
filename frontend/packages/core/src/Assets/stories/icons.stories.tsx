import * as React from "react";
import type { Meta } from "@storybook/react";

import ExperimentIcon from "../icons/ExperimentIcon";
import FireIcon from "../icons/FireIcon";
import GemIcon from "../icons/GemIcon";
import PlusIcon from "../icons/PlusIcon";
import RocketIcon from "../icons/RocketIcon";
import SirenIcon from "../icons/SirenIcon";
import SlackIcon from "../icons/SlackIcon";

export const AllIcons = () => (
  <div>
    <FireIcon />
    <PlusIcon />
    <ExperimentIcon />
    <GemIcon />
    <RocketIcon />
    <SirenIcon />
    <SlackIcon />
  </div>
);

export default {
  title: "Core/Assets/Icons",
  component: AllIcons,
} as Meta;
