import * as React from "react";
import styled from "@emotion/styled";
import {
  Grid,
  Step as MuiStep,
  StepConnector as MuiStepConnector,
  StepLabel as MuiStepLabel,
  Stepper as MuiStepper,
} from "@material-ui/core";
import MuiCheckIcon from "@material-ui/icons/Check";
import PriorityHighIcon from "@material-ui/icons/PriorityHigh";

const StepContainer = styled.div({
  ".MuiStepper-root": {
    padding: "0",
  },
  ".MuiGrid-container": {
    padding: "16px 0",
  },
});

const Circle = styled.div((props: { background: string; border: string }) => ({
  backgroundColor: props.background,
  border: props.border,
  boxSizing: "border-box",
  borderRadius: "50%",
  height: "24px",
  width: "24px",
  top: "24px",
}));

const DefaultIcon = styled.div((props: { font: string }) => ({
  height: "100%",
  width: "100%",
  color: props.font,
  display: "flex",
  alignItems: "center",
  justifyContent: "center",
  fontSize: "14px",
  fontWeight: 500,
  lineHeight: "18px",
}));

const CheckIcon = styled(MuiCheckIcon)((props: { font: string }) => ({
  fill: props.font,
  padding: "8px",
}));

const ClearIcon = CheckIcon.withComponent(PriorityHighIcon);

const StepConnector = styled(MuiStepConnector)((props: { completed?: boolean }) => ({
  ".MuiStepConnector-line": {
    height: "5px",
    background: props.completed ? "#3548D4" : "#E7E7EA",
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
    border: "1px solid #3548D4",
    font: "#3548D4",
  },
  pending: {
    background: "#E7E7EA",
    border: "#E7E7EA",
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

const StepLabel = styled(Grid)(
  {
    fontWeight: 500,
    fontSize: "14px",
  },
  props => ({
    color: props["data-active"] ? "#0D1030" : "rgba(13, 16, 48, 0.38)",
  })
);

export interface StepProps {
  label: string;
  error?: boolean;
}

const Step: React.FC<StepProps> = ({ children }) => <>{children}</>;

export interface StepperProps {
  activeStep: number;
  children?: React.ReactElement<StepProps>[] | React.ReactElement<StepProps>;
}

const Stepper = ({ activeStep, children }: StepperProps) => (
  <StepContainer>
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
    <Grid container justify="space-between">
      {React.Children.map(children, (step: any, idx: number) => (
        <StepLabel item data-active={idx === activeStep}>
          {step.props.label}
        </StepLabel>
      ))}
    </Grid>
  </StepContainer>
);

export { Stepper, Step };
