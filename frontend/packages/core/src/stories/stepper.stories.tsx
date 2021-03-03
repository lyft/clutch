import * as React from "react";
import styled from "@emotion/styled";
import type { Meta } from "@storybook/react";

import { Button, ButtonGroup } from "../button";
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
  const [curStep, setCurStep] = React.useState(activeStep || 0);

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
        <Text>
Step{curStep + 1}
{' '}
content
</Text>

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

const FailureTemplate = ({ failedStep = 2, activeStep }: StepperProps & { failedStep: number }) => {
  const [curStep, setCurStep] = React.useState(activeStep || 0);
  const stepCount = 4;

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
          <Step
            error={curStep === failedStep && index === failedStep}
            key={index} // eslint-disable-line react/no-array-index-key
            label={`Step ${index + 1}`}
          />
        ))}
      </Stepper>
      <div>
        <Text>
Step{curStep + 1}
{' '}
content
</Text>

        {curStep === stepCount - 1 || curStep === failedStep ? (
          <ButtonGroup justify="flex-start">
            <Button onClick={handleReset} text="Reset" variant="neutral" />
          </ButtonGroup>
        ) : (
          <ButtonGroup justify="flex-start">
            <Button disabled={curStep === 0} onClick={handleBack} text="Back" variant="neutral" />
            <Button onClick={handleNext} text={curStep === stepCount ? "Finish" : "Next"} />
          </ButtonGroup>
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
Failure.args = {
  failedStep: 2,
};
