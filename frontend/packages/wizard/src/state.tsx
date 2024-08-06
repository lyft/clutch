import React from "react";

enum WizardActionType {
  NEXT,
  BACK,
  RESET,
  GO_TO_STEP,
}

interface WizardAction {
  type: WizardActionType;
  step?: number;
}

interface StateProps {
  activeStep: number;
}

const reducer = (state: StateProps, action: WizardAction): StateProps => {
  switch (action.type) {
    case WizardActionType.NEXT:
      return {
        ...state,
        activeStep: state.activeStep + 1,
      };
    case WizardActionType.BACK:
      return {
        ...state,
        activeStep: state.activeStep > 0 ? state.activeStep - 1 : 0,
      };
    case WizardActionType.RESET:
      return {
        ...state,
        activeStep: 0,
      };
    case WizardActionType.GO_TO_STEP:
      return {
        ...state,
        activeStep: action.step,
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

export { WizardActionType, useWizardState };
