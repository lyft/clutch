import React from "react";

enum WizardActionType {
  NEXT,
  BACK,
  RESET,
  GO_TO_STEP,
  ADD_COMPLETED_STEP,
}

interface WizardAction {
  type: WizardActionType;
  step?: number;
}

interface StateProps {
  activeStep: number;
  nextStepToComplete: number;
  completed: {
    [key: number]: boolean;
  };
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
    case WizardActionType.ADD_COMPLETED_STEP:
      return {
        ...state,
        completed: {
          ...state.completed,
          [action.step]: true,
        },
        nextStepToComplete: state.nextStepToComplete + 1,
      };
    default:
      throw new Error(`Unknown wizard state: ${action}`);
  }
};

const initialState = {
  activeStep: 0,
  completed: {},
  nextStepToComplete: 0,
};

const useWizardState = (): [StateProps, React.Dispatch<WizardAction>] => {
  return React.useReducer(reducer, initialState);
};

export { WizardActionType, useWizardState };
