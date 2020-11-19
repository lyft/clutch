import React from "react";
import styled from "@emotion/styled";
import {
  Grid,
  Step as MuiStep,
  StepConnector as MuiStepConnector,
  StepLabel as MuiStepLabel,
  Stepper as MuiStepper,
} from "@material-ui/core";
import MuiCheckIcon from "@material-ui/icons/Check";
import MuiClearIcon from "@material-ui/icons/Clear";

import { Button } from "./button";

const Circle = styled.div((props: { background: string; border: string }) => ({
  backgroundColor: props.background,
  border: props.border,
  boxSizing: "border-box",
  borderRadius: "50%",
  height: "1.5rem",
  width: "1.5rem",
  top: "1.5rem",
}));

const DefaultIcon = styled.div((props: { font: string }) => ({
  height: "100%",
  width: "100%",
  color: props.font,
  display: "flex",
  alignItems: "center",
  textAlign: "center",
  justifyContent: "center",
  fontSize: "0.875rem",
  fontWeight: 500,
  lineHeight: "1.125rem",
}));

const CheckIcon = styled(MuiCheckIcon)({
  fill: "#FFFFFF",
  padding: "0.5rem",
});

const ClearIcon = styled(MuiClearIcon)({
  fill: "#FFFFFF",
  padding: "0.5rem",
});

const StepConnector = styled(MuiStepConnector)((props: { completed?: boolean }) => ({
  ".MuiStepConnector-line": {
    height: "0.313rem",
    background: props.completed ? "#3548D4" : "rgba(13, 16, 48, 0.12)",
    border: "0",
  },
}));

const StepLabelIcon = styled(MuiStepLabel)({
  ".MuiStepLabel-iconContainer": {
    padding: "0",
  },
});

type StepIconVariant = "active" | "pending" | "success" | "failed";
export interface StepIconProps {
  index: number;
  variant: StepIconVariant;
}

const stepIconVariants = {
  active: {
    background: "#FFFFFF",
    border: "0.063rem solid #3548D4",
    font: "#3548D4",
  },
  pending: {
    background: "rgba(13, 16, 48, 0.12)",
    border: "rgba(13, 16, 48, 0.12)",
    font: "rgba(13, 16, 48, 0.38)",
  },
  success: {
    background: "#3548D4",
    border: "#3548D4",
    font: "#3548D4",
  },
  failed: {
    background: "#DB3615",
    border: "#DB3615",
    font: "#DB3615",
  },
};

const StepIcon: React.FC<StepIconProps> = ({ index, variant }) => {
  const color = stepIconVariants[variant || "pending"];
  let Icon = <>{index}</>;
  if (variant === "success") {
    Icon = <CheckIcon fontSize="large" />;
  } else if (variant === "failed") {
    Icon = <ClearIcon fontSize="large" />;
  }
  return (
    <Circle background={color.background} border={color.border}>
      <DefaultIcon font={color.font}>{Icon}</DefaultIcon>
    </Circle>
  );
};

const StepLabel = styled(Grid)({
  fontWeight: 500,
  fontSize: "0.875rem",
  lineHeight: "1.125rem",
});

export interface StepProps {
  label: string;
  error?: boolean;
}

const Step: React.FC<StepProps> = ({ children }) => <>{children}</>;

export interface StepperProps {
  children?: React.ReactElement<StepProps>[] | React.ReactElement<StepProps>;
}

const Stepper: React.FC<StepperProps> = ({ children }) => {
  const [activeStep, setActiveStep] = React.useState(0);
  const steps = React.Children.toArray(children) as React.ReactElement<StepProps>[];
  const stepCount = steps.length - 1;

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
    <div>
      <MuiStepper activeStep={activeStep + 1} connector={<StepConnector />}>
        {React.Children.map(children, (step: any, idx: number) => {
          const stepProps = {
            index: idx + 1,
            variant: "pending" as StepIconVariant,
          };
          if (idx === activeStep) {
            stepProps.variant = "active";
          } else if (idx < activeStep) {
            stepProps.variant = step.props.error ? "failed" : "success";
          }

          return (
            <MuiStep key={step.label} style={{ padding: "0" }}>
              <StepLabelIcon icon={<StepIcon {...stepProps} />} />
            </MuiStep>
          );
        })}
      </MuiStepper>
      <Grid style={{ padding: "0 1.5rem 1.5rem 1.5rem" }} container justify="space-between">
        {React.Children.map(children, (step: any) => (
          <StepLabel item>{step.props.label}</StepLabel>
        ))}
      </Grid>

      <div>
        {steps[activeStep]}
        {activeStep === stepCount || steps[(activeStep || 1) - 1].props.error ? (
          <Button onClick={handleReset} text="Reset" variant="neutral" />
        ) : (
          <div>
            <Button
              disabled={activeStep === 0}
              onClick={handleBack}
              text="Back"
              variant="neutral"
            />
            <Button onClick={handleNext} text={activeStep === stepCount ? "Finish" : "Next"} />
          </div>
        )}
      </div>
    </div>
  );
};

export { Stepper, Step };
