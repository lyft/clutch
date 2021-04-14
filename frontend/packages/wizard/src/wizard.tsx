import React from "react";
import { useLocation, useNavigate, useSearchParams } from "react-router-dom";
import { Button, ButtonGroup, Step, Stepper, Warning, WizardContext } from "@clutch-sh/core";
import type { ManagerLayout } from "@clutch-sh/data-layout";
import { DataLayoutContext, useDataLayoutManager } from "@clutch-sh/data-layout";
import styled from "@emotion/styled";
import { Container as MuiContainer, Grid, Paper as MuiPaper, Typography } from "@material-ui/core";

import { useWizardState, WizardAction } from "./state";
import type { WizardStepProps } from "./step";

const Heading = styled(Typography)({
  paddingBottom: "16px",
  fontWeight: 700,
  fontSize: "26px",
});

interface WizardProps extends Pick<ContainerProps, "width"> {
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

interface ContainerProps {
  width?: "default" | "full";
}

const Container = styled(MuiContainer)<ContainerProps>(
  {
    padding: "32px",
    maxWidth: "unset",
  },
  props => ({
    width: props.width === "full" ? "100%" : "800px",
  })
);

const Paper = styled(MuiPaper)({
  boxShadow: "0px 5px 15px rgba(53, 72, 212, 0.2)",
  padding: "32px",
});

const Wizard = ({ heading, width = "default", dataLayout, children }: WizardProps) => {
  const [state, dispatch] = useWizardState();
  const [wizardStepData, setWizardStepData] = React.useState<WizardStepData>({});
  const [globalWarnings, setGlobalWarnings] = React.useState<string[]>([]);
  const dataLayoutManager = useDataLayoutManager(dataLayout);
  const [, setSearchParams] = useSearchParams();
  const location = useLocation();
  const navigate = useNavigate();
  const origin = location.state?.origin;

  const navigateToOrigin = () => {
    if (origin !== undefined) {
      navigate(origin);
    }
  };

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
      setHasError: (hasError: boolean) => {
        updateStepData(child.type.name, { hasError });
      },
      displayWarnings: (warnings: string[]) => {
        setGlobalWarnings(warnings);
      },
      onBack: () => {
        setGlobalWarnings([]);
        setSearchParams({});
        dispatch(WizardAction.BACK);
        navigateToOrigin();
      },
    };
  };

  const lastStepIndex = React.Children.count(children) - 1;
  // If our wizard only has 1 step, it doesn't make sense to put a restart button
  const isMultistep = lastStepIndex > 0;
  const steps = React.Children.map(children, (child: WizardChildren) => {
    const isLoading = wizardStepData[child.type.name]?.isLoading || false;
    const hasError = wizardStepData[child.type.name]?.hasError;
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
          {((state.activeStep === lastStepIndex && !isLoading && isMultistep) || hasError) && (
            <ButtonGroup>
              <Button
                text="Start Over"
                onClick={() => {
                  dataLayoutManager.reset();
                  setSearchParams({});
                  dispatch(WizardAction.RESET);
                  navigateToOrigin();
                }}
              />
            </ButtonGroup>
          )}
        </Grid>
      </>
    );
  });

  const removeWarning = (warning: string) => {
    setGlobalWarnings(globalWarnings.filter(w => w !== warning));
  };

  return (
    <Container width={width}>
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
              const { name } = child.props;
              const hasError = wizardStepData[child.type.name]?.hasError;
              return <Step key={name} label={name} error={hasError} />;
            })}
          </Stepper>
          <Paper elevation={0}>{steps[state.activeStep]}</Paper>
        </Grid>
      </Grid>
      {globalWarnings.map(error => (
        <Warning key={error} message={error} onClose={() => removeWarning(error)} />
      ))}
    </Container>
  );
};

export default Wizard;
