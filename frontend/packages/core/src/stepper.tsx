import * as React from "react";
import styled from "@emotion/styled";
import {
  Step as MuiStep,
  StepConnector as MuiStepConnector,
  StepLabel as MuiStepLabel,
  Stepper as MuiStepper,
} from "@material-ui/core";
import MuiCheckIcon from "@material-ui/icons/Check";
import PriorityHighIcon from "@material-ui/icons/PriorityHigh";

const StepContainer = styled.div({
  margin: "0px 2px 30px 2px",
  ".MuiStepLabel-label": {
    fontWeight: 500,
    fontSize: "14px",
    color: "rgba(13, 16, 48, 0.38)",
  },
  ".MuiStepLabel-label.MuiStepLabel-active": {
    color: "#0d1030",
  },
  ".MuiStepper-root": {
    background: "transparent",
    padding: "0",
  },
  ".MuiGrid-container": {
    padding: "16px 0",
  },
  ".MuiStepLabel-labelContainer": {
    width: "unset",
  },
  ".MuiStepConnector-alternativeLabel": {
    top: "10px",
    right: "calc(50% + 8px)",
    left: "calc(-50% + 8px)",
    zIndex: 10,
  },
  ".MuiStepLabel-iconContainer": {
    zIndex: 20,
  },
  ".MuiStep-root": {
    padding: "0",
  },
  ".MuiStep-root:first-of-type": {
    ".MuiStepLabel-root": {
      alignItems: "flex-start",
    },
  },
  ".MuiStep-root:nth-of-type(2)": {
    ".MuiStepConnector-alternativeLabel": {
      left: "calc(-100%)",
    },
  },
  ".MuiStep-root:last-of-type": {
    ".MuiStepLabel-root": {
      alignItems: "flex-end",
    },

    ".MuiStepConnector-alternativeLabel": {
      right: "0px",
    },
  },

  ".MuiStepConnector-line": {
    height: "5px",
    border: 0,
    backgroundColor: "#E7E7EA",
    borderRadius: "4px",
  },

  ".MuiStepConnector-active .MuiStepConnector-line": {
    backgroundColor: "#3548D4",
  },

  ".MuiStepConnector-completed .MuiStepConnector-line": {
    backgroundColor: "#3548D4",
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
    <MuiStepper activeStep={activeStep} connector={<MuiStepConnector />} alternativeLabel>
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
          <MuiStep key={step.props.label}>
            <MuiStepLabel StepIconComponent={() => <StepIcon {...stepProps} />}>
              {step.props.label ?? `Step ${idx + 1}`}
            </MuiStepLabel>
          </MuiStep>
        );
      })}
    </MuiStepper>
  </StepContainer>
);

export { Stepper, Step };
