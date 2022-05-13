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
    <FireIcon size="xsmall" />
    <FireIcon size="small" />
    <FireIcon size="medium" />
    <FireIcon size="large" />
    <PlusIcon disabled />
    <PlusIcon />
    <ExperimentIcon size="xsmall" />
    <ExperimentIcon size="small" />
    <ExperimentIcon size="medium" />
    <ExperimentIcon size="large" />
    <RocketIcon size="xsmall" />
    <RocketIcon size="small" />
    <RocketIcon size="medium" />
    <RocketIcon size="large" />
    <SirenIcon size="xsmall" />
    <SirenIcon size="small" />
    <SirenIcon size="medium" />
    <SirenIcon size="large" />
    <GemIcon size="xsmall" />
    <GemIcon size="small" />
    <GemIcon size="medium" />
    <GemIcon size="large" />
    <SlackIcon size="xsmall" />
    <SlackIcon size="small" />
    <SlackIcon size="medium" />
    <SlackIcon size="large" />
  </div>
);

export default {
  title: "Core/Assets/Icons",
  component: AllIcons,
} as Meta;
