import * as React from "react";
import type { Meta } from "@storybook/react";

import { Button, ButtonGroup } from "../button";
import Grid from "../grid";
import type { StepperProps } from "../stepper";
import { Step, Stepper } from "../stepper";
import { Typography } from "../typography";

export default {
  title: "Core/Stepper",
  component: Stepper,
  argTypes: {
    orientation: {
      options: ["horizontal", "vertical"],
      defaultValue: "horizontal",
      control: {
        type: "select",
      },
    },
  },
} as Meta;

const PrimaryTemplate = ({
  stepCount,
  activeStep,
  orientation,
}: StepperProps & { stepCount: number }) => {
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
      <Grid container direction={orientation === "horizontal" ? "column" : "row"}>
        <Grid item xs={1}>
          <Stepper activeStep={curStep} orientation={orientation}>
            {Array(stepCount)
              .fill(null)
              .map((_, index: number) => (
                // eslint-disable-next-line react/no-array-index-key
                <Step key={index} label={`Step ${index + 1}`} />
              ))}
          </Stepper>
        </Grid>
        <Grid container item xs={11} alignContent="center" justifyContent="center">
          <Typography variant="body3">Step{curStep + 1} content</Typography>
        </Grid>
      </Grid>
      <ButtonGroup justify="flex-start">
        {curStep === stepCount - 1 ? (
          <Button onClick={handleReset} text="Reset" variant="neutral" />
        ) : (
          <>
            <Button disabled={curStep === 0} onClick={handleBack} text="Back" variant="neutral" />
            <Button onClick={handleNext} text={curStep === stepCount ? "Finish" : "Next"} />
          </>
        )}
      </ButtonGroup>
    </>
  );
};

const FailureTemplate = ({
  failedStep = 2,
  activeStep,
  orientation,
}: StepperProps & { failedStep: number }) => {
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
      <Grid container direction={orientation === "horizontal" ? "column" : "row"}>
        <Grid item xs={1}>
          <Stepper activeStep={curStep} orientation={orientation}>
            {Array(stepCount)
              .fill(null)
              .map((_, index: number) => (
                <Step
                  error={curStep === failedStep && index === failedStep}
                  key={index} // eslint-disable-line react/no-array-index-key
                  label={`Step ${index + 1}`}
                />
              ))}
          </Stepper>
        </Grid>
        <Grid container item xs={11} alignContent="center" justifyContent="center">
          <Typography variant="body3">Step{curStep + 1} content</Typography>
        </Grid>
      </Grid>
      <ButtonGroup justify="flex-start">
        {curStep === stepCount - 1 || curStep === failedStep ? (
          <Button onClick={handleReset} text="Reset" variant="neutral" />
        ) : (
          <>
            <Button disabled={curStep === 0} onClick={handleBack} text="Back" variant="neutral" />
            <Button onClick={handleNext} text={curStep === stepCount ? "Finish" : "Next"} />
          </>
        )}
      </ButtonGroup>
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
