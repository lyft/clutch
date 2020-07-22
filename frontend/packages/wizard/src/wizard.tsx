import React from "react";
import { Button, Warning, WizardContext } from "@clutch-sh/core";
import type { ManagerLayout } from "@clutch-sh/data-layout";
import { DataLayoutContext, useDataLayoutManager } from "@clutch-sh/data-layout";
import {
  Grid,
  Step,
  StepConnector as MUIStepConnector,
  StepContent,
  StepLabel,
  Stepper,
  Typography,
} from "@material-ui/core";
import { makeStyles, withStyles } from "@material-ui/core/styles";
import Check from "@material-ui/icons/Check";
import ErrorOutlineIcon from "@material-ui/icons/ErrorOutline";
import FirstPageIcon from "@material-ui/icons/FirstPage";
import clsx from "clsx";
import styled from "styled-components";

import { useWizardState, WizardAction } from "./state";
import type { WizardStepProps } from "./step";

const Heading = styled(Typography)`
  padding-left: 2.5%;
`;

const StepConnector = withStyles({
  alternativeLabel: {
    top: 10,
    left: "calc(-50% + 16px)",
    right: "calc(50% + 16px)",
  },
  active: {
    "& $line": {
      borderColor: "#02acbe",
    },
  },
  completed: {
    "& $line": {
      borderColor: "#02acbe",
    },
  },
  line: {
    borderColor: "#eaeaf0",
    borderTopWidth: 3,
    borderRadius: 1,
  },
})(MUIStepConnector);

const useQontoStepIconStyles = makeStyles({
  root: {
    color: "#eaeaf0",
    display: "flex",
    height: 22,
    alignItems: "center",
  },
  active: {
    color: "#02acbe",
  },
});

const CircleIcon = styled.div`
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background-color: currentColor;
`;

const CheckmarkIcon = styled(Check)`
  ${({ theme }) => `
  color: ${theme.palette.accent.main};
  z-index: 1;
  font-size: 18px;
  `}
`;

interface StepIconProps {
  active: boolean;
  completed: boolean;
  error: boolean;
}

const StepIcon: React.FC<StepIconProps> = ({ active, completed, error }) => {
  const classes = useQontoStepIconStyles();

  return (
    <div
      className={clsx(classes.root, {
        [classes.active]: active,
      })}
    >
      {completed ? <CheckmarkIcon /> : error ? <ErrorOutlineIcon /> : <CircleIcon />}
    </div>
  );
};

const StyledStepper = styled(Stepper)`
  background-color: transparent;
`;

interface SpacerProps {
  margin?: string;
}

const Spacer = styled.div<SpacerProps>`
  ${({ margin }) => `
  margin: ${Number(margin || 1) * 10}px;
  `}
`;

interface WizardProps {
  heading?: string;
  dataLayout: ManagerLayout;
}

export interface WizardChild {
  name: string;
}

interface WizardChildren extends JSX.Element {
  value: WizardStepProps;
}

interface WizardStepData {
  [index: string]: any;
}

const SizedContainer = styled(Grid)`
  display: inline;
`;

const StartOverIcon = styled(FirstPageIcon)`
  transform: rotate(90deg);
`;

const Wizard: React.FC<WizardProps> = ({ heading, dataLayout, children }) => {
  const [state, dispatch] = useWizardState();
  const [wizardStepData, setWizardStepData] = React.useState<WizardStepData>({});
  const [gloablWarnings, setGlobalWarnings] = React.useState<string[]>([]);
  const dataLayoutManager = useDataLayoutManager(dataLayout);

  const updateStepData = (stepName: string, data: object) => {
    const updatedData = {
      ...setWizardStepData?.[stepName],
      ...data,
    };
    const stepData = { [stepName]: updatedData };
    setWizardStepData({ ...wizardStepData, ...stepData });
  };

  const handleNext = () => {
    dispatch(WizardAction.NEXT);
  };

  const context = (child: JSX.Element) => {
    return {
      onSubmit: wizardStepData?.[child.type.name]?.onSubmit || handleNext,
      setOnSubmit: (f: (...args: any[]) => void) => {
        updateStepData(child.type.name, { onSubmit: f(handleNext) });
      },
      setIsLoading: (isLoading: boolean) => {
        updateStepData(child.type.name, { isLoading });
      },
      displayWarnings: (warnings: string[]) => {
        setGlobalWarnings(warnings);
      },
      onBack: () => {
        setGlobalWarnings([]);
        dispatch(WizardAction.BACK);
      },
    };
  };

  const lastStepIndex = React.Children.count(children) - 1;
  // If our wizard only has 1 step, it doesn't make sense to put a restart button
  const isMultistep = lastStepIndex > 0;
  const steps = React.Children.map(children, (child: WizardChildren, idx: number) => {
    const hasError = (wizardStepData[child.type.name]?.errors?.length || 0) !== 0;
    const isLoading = wizardStepData[child.type.name]?.isLoading || false;
    return (
      <Step key={child.props.name} expanded={idx <= state.activeStep}>
        <StepLabel StepIconComponent={StepIcon} error={hasError}>
          {child.props.name}
        </StepLabel>
        <StepContent
          TransitionProps={{ appear: true }}
          style={{ display: idx === state.activeStep ? "block" : "none" }}
        >
          <DataLayoutContext.Provider value={dataLayoutManager}>
            <WizardContext.Provider value={() => context(child)}>
              <Grid container direction="column" justify="center" alignItems="center">
                {child}
              </Grid>
            </WizardContext.Provider>
          </DataLayoutContext.Provider>
          <Grid container justify="center">
            {state.activeStep === lastStepIndex && !isLoading && isMultistep && (
              <Button
                onClick={() => dispatch(WizardAction.RESET)}
                text="Start Over"
                endIcon={<StartOverIcon />}
              />
            )}
          </Grid>
        </StepContent>
      </Step>
    );
  });

  return (
    <Spacer margin="3">
      {heading && <Heading variant="h5">{heading}</Heading>}
      <SizedContainer
        container
        direction="column"
        justify="center"
        alignItems="stretch"
        style={{ display: "inline" }}
      >
        <Grid item>
          <StyledStepper
            orientation="vertical"
            activeStep={state.activeStep}
            connector={<StepConnector />}
          >
            {steps}
          </StyledStepper>
        </Grid>
      </SizedContainer>
      {gloablWarnings.map(error => (
        <Warning key={error} message={error} />
      ))}
    </Spacer>
  );
};

export default Wizard;
