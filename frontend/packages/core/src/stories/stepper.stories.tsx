import React from "react";
import styled from "@emotion/styled";
import type { Meta } from "@storybook/react";

import { Step, Stepper } from "../stepper";

const Text = styled.div({
  textAlign: "center",
});

export default {
  title: "Core/Stepper",
  component: Stepper,
} as Meta;

const PrimaryTemplate = ({ stepCount }: { stepCount: number }) => (
  <Stepper>
    {[...Array(stepCount)].map((_, index: number) => (
      // eslint-disable-next-line react/no-array-index-key
      <Step key={index} label={`Step ${index + 1}`}>
        <Text>Step {index + 1} content</Text>
      </Step>
    ))}
  </Stepper>
);

const FailureTemplate = () => (
  <Stepper>
    <Step label="Step 1">
      <Text>First step content</Text>
    </Step>
    <Step label="Step 2" error>
      <Text>Second step content</Text>
    </Step>
    <Step label="Step 3">
      <Text>Third step content</Text>
    </Step>
    <Step label="Step 4">
      <Text>Fourth step content</Text>
    </Step>
  </Stepper>
);

export const Primary = PrimaryTemplate.bind({});
Primary.args = {
  stepCount: 3,
};

export const Failure = FailureTemplate.bind({});
