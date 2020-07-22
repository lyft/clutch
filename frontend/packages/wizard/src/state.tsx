import React from "react";

enum WizardAction {
  NEXT,
  BACK,
  RESET,
}

interface StateProps {
  activeStep: number;
}

const reducer = (state: StateProps, action: WizardAction): StateProps => {
  switch (action) {
    case WizardAction.NEXT:
      return {
        ...state,
        activeStep: state.activeStep + 1,
      };
    case WizardAction.BACK:
      return {
        ...state,
        activeStep: state.activeStep - 1,
      };
    case WizardAction.RESET:
      return {
        ...state,
        activeStep: 0,
      };
    default:
      throw new Error(`Unknown wizard state: ${action}`);
  }
};

const initialState = {
  activeStep: 0,
};

const useWizardState = (): [StateProps, React.Dispatch<WizardAction>] => {
  return React.useReducer(reducer, initialState);
};

export { WizardAction, useWizardState };
