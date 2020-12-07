import * as React from "react";
import styled from "@emotion/styled";
import type { Meta } from "@storybook/react";

import { Button } from "../button";
import type { StepperProps } from "../stepper";
import { Step, Stepper } from "../stepper";

const Text = styled.div({
  textAlign: "center",
});

export default {
  title: "Core/Stepper",
  component: Stepper,
} as Meta;

const PrimaryTemplate = ({ stepCount, activeStep }: StepperProps & { stepCount: number }) => {
  const [curStep, setCurStep] = React.useState(activeStep);

  const handleNext = () => {
    setCurStep(prevActiveStep => prevActiveStep + 1);
  };

  const handleBack = () => {
    setCurStep(prevActiveStep => prevActiveStep - 1);
  };

  const handleReset = () => {
    setCurStep(0);
  };

  return (
    <>
      <Stepper activeStep={curStep}>
        {[...Array(stepCount)].map((_, index: number) => (
          // eslint-disable-next-line react/no-array-index-key
          <Step key={index} label={`Step ${index + 1}`} />
        ))}
      </Stepper>
      <div>
        <Text>Step {curStep + 1} content</Text>

        {curStep === stepCount - 1 ? (
          <Button onClick={handleReset} text="Reset" variant="neutral" />
        ) : (
          <div>
            <Button disabled={curStep === 0} onClick={handleBack} text="Back" variant="neutral" />
            <Button onClick={handleNext} text={curStep === stepCount ? "Finish" : "Next"} />
          </div>
        )}
      </div>
    </>
  );
};

const FailureTemplate = () => {
  const [activeStep, setActiveStep] = React.useState(0);

  const handleNext = () => {
    setActiveStep(prevActiveStep => prevActiveStep + 1);
  };

  const handleBack = () => {
    setActiveStep(prevActiveStep => prevActiveStep - 1);
  };

  const handleReset = () => {
    setActiveStep(0);
  };

  return (
    <>
      <Stepper activeStep={activeStep}>
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
      <div>
        <Text>Step {activeStep + 1} content</Text>

        {activeStep === 2 ? (
          <Button onClick={handleReset} text="Reset" variant="neutral" />
        ) : (
          <div>
            <Button
              disabled={activeStep === 0}
              onClick={handleBack}
              text="Back"
              variant="neutral"
            />
            <Button onClick={handleNext} text={activeStep === 3 ? "Finish" : "Next"} />
          </div>
        )}
      </div>
    </>
  );
};

export const Primary = PrimaryTemplate.bind({});
Primary.args = {
  stepCount: 3,
};

export const Failure = FailureTemplate.bind({});
