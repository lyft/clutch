import React from "react";
import { Warning, Step, Stepper, WizardContext, ButtonGroup } from "@clutch-sh/core";
import type { ManagerLayout } from "@clutch-sh/data-layout";
import { DataLayoutContext, useDataLayoutManager } from "@clutch-sh/data-layout";
import {
  Container as MuiContainer,
  Grid,
  Typography,
} from "@material-ui/core";
import styled from "@emotion/styled";

import { useWizardState, WizardAction } from "./state";
import type { WizardStepProps } from "./step";

const Heading = styled(Typography)({
  paddingBottom: "16px",
  fontWeight: 700,
  fontSize: "26px",
});

interface WizardProps {
  heading?: string;
  dataLayout: ManagerLayout;
  children: React.ReactElement<WizardStepProps> | React.ReactElement<WizardStepProps>[];
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

const Container = styled(MuiContainer)({
  padding: "32px",
  maxWidth: "800px",
});

const Wizard = ({ heading, dataLayout, children }: WizardProps) => {
  const [state, dispatch] = useWizardState();
  const [wizardStepData, setWizardStepData] = React.useState<WizardStepData>({});
  const [globalWarnings, setGlobalWarnings] = React.useState<string[]>([]);
  const dataLayoutManager = useDataLayoutManager(dataLayout);

  const updateStepData = (stepName: string, data: object) => {
    setWizardStepData(prevState => {
      const updatedData = {
        ...(prevState?.[stepName] || {}),
        ...data,
      };
      const stepData = { [stepName]: updatedData };
      return { ...prevState, ...stepData };
    });
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
    const isLoading = wizardStepData[child.type.name]?.isLoading || false;
    return (
      <>
        <DataLayoutContext.Provider value={dataLayoutManager}>
          <WizardContext.Provider value={() => context(child)}>
            <Grid container direction="column" justify="center" alignItems="center">
              {child}
            </Grid>
          </WizardContext.Provider>
        </DataLayoutContext.Provider>
        <Grid container justify="center">
          {state.activeStep === lastStepIndex && !isLoading && isMultistep && (
          <ButtonGroup
            justify="flex-end"
            buttons={[
              {
                text: "Start Over",
                onClick: () => dispatch(WizardAction.RESET),
              },
            ]}
          />
          )}
        </Grid>
      </>
    );
  });

  const removeWarning = (warning: string) => {
    setGlobalWarnings(globalWarnings.filter(w => w !== warning));
  };

  return (
    <Container>
      <Grid
        container
        direction="column"
        justify="center"
        alignItems="stretch"
        style={{ display: "inline" }}
      >
      {heading && <Heading>{heading}</Heading>}
        <Grid item>
          <Stepper activeStep={state.activeStep}>
            {React.Children.map(children, (child: WizardChildren) => {
              const hasError = (wizardStepData[child.type.name]?.errors?.length || 0) !== 0;
              return <Step key={child.props.name} label={child.props.name} error={hasError} />;
            })}
          </Stepper>
          {steps[state.activeStep]}
        </Grid>
      </Grid>
      {globalWarnings.map(error => (
        <Warning key={error} message={error} onClose={() => removeWarning(error)} />
      ))}
    </Container>
  );
};

export default Wizard;
