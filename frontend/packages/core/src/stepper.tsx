import * as React from "react";
import MuiCheckIcon from "@mui/icons-material/Check";
import PriorityHighIcon from "@mui/icons-material/PriorityHigh";
import type {
  Orientation as StepperOrientation,
  StepperProps as MuiStepperProps,
} from "@mui/material";
import {
  alpha,
  Step as MuiStep,
  StepButton as MuiStepButton,
  StepConnector as MuiStepConnector,
  StepLabel as MuiStepLabel,
  Stepper as MuiStepper,
  Theme,
  useTheme,
} from "@mui/material";

import styled from "./styled";

const StepContainer = styled("div")<{ $orientation: StepperOrientation; $nonLinear: boolean }>(
  {
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
    ".MuiStep-root:nth-of-type(2)": {
      ".MuiStepConnector-alternativeLabel": {
        left: "calc(-100%)",
      },
    },
    ".MuiStep-root:last-of-type": {
      ".MuiStepConnector-alternativeLabel": {
        right: "0px",
      },
    },
  },
  props => ({ theme }: { theme: Theme }) => ({
    ".MuiStep-root:hover": {
      ".MuiStepLabel-iconContainer .icon-circle-nonlinear-pending": {
        border: props.$nonLinear ? "2px solid black" : "",
      },
      ".MuiStepLabel-label": {
        color: props.$nonLinear ? theme.palette.secondary[600] : "",
      },
    },
    ".MuiStepLabel-label": {
      fontWeight: 500,
      fontSize: "14px",
      color: theme.palette.secondary[300],
    },
    ".MuiStepLabel-label.Mui-active": {
      color: theme.palette.primary[600],
    },
    ".MuiStepLabel-label.Mui-completed": {
      color: theme.palette.secondary[600],
    },
    ...(props.$orientation === "horizontal"
      ? {
          margin: "0px 2px 30px 2px",
          ".MuiStep-root:first-of-type": {
            ".MuiStepLabel-root": {
              alignItems: "flex-start",
            },
          },
          ".MuiStep-root:last-of-type": {
            ".MuiStepLabel-root": {
              alignItems: "flex-end",
            },
          },
          ".MuiStepConnector-line": {
            height: props.$nonLinear ? "3px" : "5px",
            border: 0,
            backgroundColor: props.$nonLinear
              ? theme.palette.primary[600]
              : theme.palette.secondary[200],
            borderRadius: "4px",
          },
          ".Mui-active .MuiStepConnector-line": {
            backgroundColor: theme.palette.primary[600],
          },
          ".Mui-completed .MuiStepConnector-line": {
            backgroundColor: theme.palette.primary[600],
          },
        }
      : {
          margin: "0px 2px 8px 2px",
          ".MuiStepConnector-line": {
            borderColor: props.$nonLinear
              ? theme.palette.primary[600]
              : theme.palette.secondary[300],
          },
          ".Mui-active .MuiStepConnector-line": {
            borderColor: theme.palette.primary[600],
          },
          ".Mui-completed .MuiStepConnector-line": {
            borderColor: theme.palette.primary[600],
          },
        }),
  })
);

const Circle = styled("div")((props: { background: string; border: string }) => ({
  backgroundColor: props.background,
  border: props.border,
  boxSizing: "border-box",
  borderRadius: "50%",
  height: "24px",
  width: "24px",
  top: "24px",
}));

const DefaultIcon = styled("div")((props: { font: string }) => ({
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
  nonLinear: boolean;
}

const StepIcon: React.FC<StepIconProps> = ({ index, variant, nonLinear }) => {
  const theme = useTheme();
  const stepIconVariants = {
    active: {
      background: theme.palette.contrastColor,
      border: `2px solid ${theme.palette.primary[600]}`,
      font: theme.palette.primary[600],
    },
    pending: {
      background: nonLinear ? theme.palette.secondary[100] : theme.palette.secondary[200],
      border: nonLinear ? theme.palette.secondary[100] : theme.palette.secondary[200],
      font: nonLinear ? theme.palette.secondary[600] : alpha(theme.palette.secondary[900], 0.38),
    },
    success: {
      background: theme.palette.primary[600],
      border: theme.palette.primary[600],
      font: theme.palette.contrastColor,
    },
    failed: {
      background: theme.palette.error[600],
      border: theme.palette.error[600],
      font: theme.palette.contrastColor,
    },
  };
  const color = stepIconVariants[variant || "pending"];
  let Icon = <>{index}</>;
  if (variant === "success") {
    Icon = <CheckIcon font={color.font} fontSize="large" />;
  } else if (variant === "failed") {
    Icon = <ClearIcon font={color.font} fontSize="large" />;
  }
  const nonLinearPendingClass =
    nonLinear && variant === "pending" ? "icon-circle-nonlinear-pending" : "";
  return (
    <Circle background={color.background} border={color.border} className={nonLinearPendingClass}>
      <DefaultIcon font={color.font}>{Icon}</DefaultIcon>
    </Circle>
  );
};

/* Because these props are just used on the children of Step, they are throwing an error as unused */
/* eslint-disable react/no-unused-prop-types */
export interface StepProps {
  label: string;
  error?: boolean;
  isComplete?: boolean;
}
/* eslint-enable react/no-unused-prop-types */

const Step: React.FC<StepProps> = ({ children }) => <>{children}</>;

export interface StepperProps
  extends Pick<MuiStepperProps, "orientation">,
    Pick<MuiStepperProps, "nonLinear"> {
  activeStep: number;
  children?: React.ReactElement<StepProps>[] | React.ReactElement<StepProps>;
  handleStepClick?: (index: number) => void;
}

const Stepper = ({
  activeStep,
  orientation = "horizontal",
  children,
  nonLinear,
  handleStepClick,
}: StepperProps) => (
  <StepContainer $orientation={orientation} $nonLinear={nonLinear}>
    <MuiStepper
      nonLinear={nonLinear}
      activeStep={activeStep}
      connector={<MuiStepConnector />}
      alternativeLabel={orientation === "horizontal"}
      orientation={orientation}
    >
      {React.Children.map(children, (step: any, idx: number) => {
        const stepProps = {
          index: idx + 1,
          variant: "pending" as StepIconVariant,
          nonLinear,
        };
        const { isComplete = false } = step.props;

        if (isComplete || (!nonLinear && idx < activeStep)) {
          stepProps.variant = "success";
        } else if (idx === activeStep) {
          stepProps.variant = step.props.error ? "failed" : "active";
        }

        const label = step.props.label ?? `Step ${idx + 1}`;
        const icon = <StepIcon {...stepProps} />;

        return (
          <MuiStep key={label} completed={isComplete}>
            {nonLinear ? (
              <MuiStepButton onClick={() => handleStepClick(idx)} icon={icon}>
                {label}
              </MuiStepButton>
            ) : (
              <MuiStepLabel icon={icon}>{label}</MuiStepLabel>
            )}
          </MuiStep>
        );
      })}
    </MuiStepper>
  </StepContainer>
);

export { Stepper, Step };
