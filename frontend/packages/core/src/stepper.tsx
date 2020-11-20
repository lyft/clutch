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
  justifyContent: "center",
  fontSize: "0.875rem",
  fontWeight: 500,
  lineHeight: "1.125rem",
}));

const CheckIcon = styled(MuiCheckIcon)((props: { font: string }) => ({
  fill: props.font,
  padding: "0.5rem",
}));

const ClearIcon = CheckIcon.withComponent(MuiClearIcon);

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
    font: "#FFFFFF",
  },
  failed: {
    background: "#DB3615",
    border: "#DB3615",
    font: "#FFFFFF",
  },
};

const StepIcon: React.FC<StepIconProps> = ({ index, variant }) => {
  const color = stepIconVariants[variant || "pending"];
  let Icon = <>{index}</>;
  if (variant === "success") {
    Icon = <CheckIcon font={color.font} fontSize="large" />;
  } else if (variant === "failed") {
    Icon = <ClearIcon font={color.font} fontSize="large" />;
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
  activeStep: number;
  children?: React.ReactElement<StepProps>[] | React.ReactElement<StepProps>;
}

const Stepper: React.FC<StepperProps> = ({ activeStep, children }) => {
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
    </div>
  );
};

export { Stepper, Step };
