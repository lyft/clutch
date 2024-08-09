import React from "react";
import type { WizardNavigationProps } from "@clutch-sh/core";
import {
  Button,
  ButtonGroup,
  FeatureOn,
  Grid,
  NPSWizard,
  Paper,
  SimpleFeatureFlag,
  Step,
  Stepper,
  styled,
  Toast,
  Typography,
  useLocation,
  useNavigate,
  useSearchParams,
  WizardContext,
} from "@clutch-sh/core";
import type { ManagerLayout } from "@clutch-sh/data-layout";
import { DataLayoutContext, useDataLayoutManager } from "@clutch-sh/data-layout";
import type {
  ContainerProps as MuiContainerProps,
  StepperProps as MuiStepperProps,
} from "@mui/material";
import { alpha, Container as MuiContainer, Theme } from "@mui/material";

import { useWizardState, WizardActionType } from "./state";
import type { WizardStepProps } from "./step";

export interface WizardProps
  extends Pick<ContainerProps, "width" | "className">,
    Pick<MuiStepperProps, "orientation" | "nonLinear"> {
  children:
    | React.ReactNode
    | React.ReactElement<WizardStepProps>
    | React.ReactElement<WizardStepProps>[];
  dataLayout: ManagerLayout;
  heading?: string | React.ReactElement;
}

export interface WizardChild {
  name: string;
  showNPS?: boolean;
  confirm?: {
    startOver: boolean;
    startOverText?: string;
  };
}

interface WizardChildren extends JSX.Element {
  value: WizardStepProps;
}

interface WizardStepData {
  [index: string]: any;
}

interface ContainerProps extends Pick<MuiContainerProps, "className"> {
  width?: "default" | "full";
}

const Header = styled(Grid)<{ $orientation: MuiStepperProps["orientation"] }>(
  {
    paddingBottom: "16px",
  },
  props => ({
    marginLeft: props.$orientation === "vertical" ? "-16px" : "0",
  })
);

const Container = styled(MuiContainer)<{ $width: ContainerProps["width"] }>(
  {
    paddingBlock: "24px 32px",
    height: "calc(100% - 56px)",
  },
  props => ({
    width: props.$width === "full" ? "100%" : "800px",
  })
);

const StepperContainer = styled(Grid)<{ $orientation: MuiStepperProps["orientation"] }>(
  {},
  props => ({ theme }: { theme: Theme }) => ({
    ...(props.$orientation === "vertical" && {
      padding: "16px",
      scrollbarGutter: "stable",
      height: "fit-content",
      maxHeight: "100%",
      overflowY: "scroll",
      background: alpha(theme.palette.secondary[200], 0.35),
    }),
  })
);

const MaxHeightGrid = styled(Grid)({
  height: "100%",
});

const StyledStepContainer = styled(Grid)({
  marginTop: "-16px",
});

const StyledPaper = styled(Paper)(({ theme }: { theme: Theme }) => ({
  boxShadow: `0px 5px 15px ${alpha(theme.palette.primary[600], 0.2)}`,
  padding: "32px",
  maxHeight: "100%",
  overflowY: "scroll",
}));

const Wizard = ({
  heading,
  width = "default",
  dataLayout,
  orientation = "horizontal",
  children,
  className,
  nonLinear = false,
}: WizardProps) => {
  const [state, dispatch] = useWizardState();
  const [wizardStepData, setWizardStepData] = React.useState<WizardStepData>({});
  const [globalWarnings, setGlobalWarnings] = React.useState<string[]>([]);
  const dataLayoutManager = useDataLayoutManager(dataLayout);
  const [, setSearchParams] = useSearchParams();
  const locationState = useLocation().state as { origin?: string };
  const navigate = useNavigate();
  const [origin] = React.useState(locationState?.origin);

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
    dispatch({ type: WizardActionType.NEXT });
  };

  const handleStepClick = (step: number) => {
    dispatch({ type: WizardActionType.GO_TO_STEP, step });
  };

  const handleNavigation = (params: WizardNavigationProps, actionType: WizardActionType) => {
    setGlobalWarnings([]);
    if (!params?.keepSearch) {
      setSearchParams({});
    }
    if (params?.toOrigin && origin) {
      navigate(origin);
    } else {
      dispatch({ type: actionType });
    }
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
      onBack: (params: WizardNavigationProps) => {
        handleNavigation(params, WizardActionType.BACK);
      },
      onNext: (params: WizardNavigationProps) => {
        handleNavigation(params, WizardActionType.NEXT);
      },
    };
  };

  // toArray will exclude any null Children.
  const filteredChildren = React.Children.toArray(children);

  const lastStepIndex = filteredChildren.length - 1;
  // If our wizard only has 1 step, it doesn't make sense to put a restart button
  const isMultistep = lastStepIndex > 0;
  const steps = filteredChildren.map((child: WizardChildren) => {
    const isLoading = wizardStepData[child.type.name]?.isLoading || false;
    const hasError = wizardStepData[child.type.name]?.hasError;
    const {
      props: {
        showNPS = true,
        confirm = {
          startOver: true,
          startOverText: "Start Over",
        },
      },
    } = child;

    return (
      <>
        <DataLayoutContext.Provider value={dataLayoutManager}>
          {/* eslint-disable-next-line react/jsx-no-constructed-context-values */}
          <WizardContext.Provider value={() => context(child)}>
            <Grid container direction="column" justifyContent="center" alignItems="center">
              {child}
            </Grid>
          </WizardContext.Provider>
        </DataLayoutContext.Provider>
        <Grid container justifyContent="center">
          {((state.activeStep === lastStepIndex && !isLoading) || hasError) && (
            <>
              {showNPS && (
                <SimpleFeatureFlag feature="npsWizard">
                  <FeatureOn>
                    <NPSWizard />
                  </FeatureOn>
                </SimpleFeatureFlag>
              )}
              {(isMultistep || hasError) && confirm.startOver && (
                <ButtonGroup>
                  <Button
                    text={confirm.startOverText ?? "Start Over"}
                    onClick={() => {
                      dataLayoutManager.reset();
                      setSearchParams({});
                      dispatch({ type: WizardActionType.RESET });
                      if (origin) {
                        navigate(origin);
                      }
                    }}
                  />
                </ButtonGroup>
              )}
            </>
          )}
        </Grid>
      </>
    );
  });

  const removeWarning = (warning: string) => {
    setGlobalWarnings(globalWarnings.filter(w => w !== warning));
  };

  return (
    <Container $width={width} maxWidth={false} className={className}>
      <MaxHeightGrid container alignItems="stretch">
        {heading && (
          <Header item $orientation={orientation}>
            {React.isValidElement(heading) ? (
              heading
            ) : (
              <Typography variant="h2">{heading}</Typography>
            )}
          </Header>
        )}
        <MaxHeightGrid
          container
          item
          direction={orientation === "vertical" ? "row" : "column"}
          wrap="nowrap"
          spacing={2}
          marginTop={0}
        >
          <StepperContainer item xs="auto" $orientation={orientation}>
            <Stepper
              activeStep={state.activeStep}
              orientation={orientation}
              nonLinear={nonLinear}
              handleStepClick={handleStepClick}
            >
              {filteredChildren.map((child: WizardChildren) => {
                const { name, isComplete } = child.props;
                const hasError = wizardStepData[child.type.name]?.hasError;
                return <Step key={name} label={name} error={hasError} isComplete={isComplete} />;
              })}
            </Stepper>
          </StepperContainer>
          <StyledStepContainer item xs={12}>
            <StyledPaper elevation={0}>{steps[state.activeStep]}</StyledPaper>
          </StyledStepContainer>
        </MaxHeightGrid>
      </MaxHeightGrid>
      {globalWarnings.map(error => (
        <Toast key={error} onClose={() => removeWarning(error)}>
          {error}
        </Toast>
      ))}
    </Container>
  );
};

export default Wizard;
