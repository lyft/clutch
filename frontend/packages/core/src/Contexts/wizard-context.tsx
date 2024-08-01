import React from "react";
import type { NavigationProps } from "@clutch-sh/wizard";

export interface ContextProps {
  displayWarnings: (warnings: string[]) => void;
  onBack: (params?: NavigationProps) => void;
  onNext: (params?: NavigationProps) => void;
  onSubmit: () => void;
  setOnSubmit: (f: (...args: any[]) => void) => void;
  setIsLoading: (isLoading: boolean) => void;
  setHasError: (hasError: boolean) => void;
  getNextStepToComplete: () => number;
  onComplete: (id: number) => void;
}

const WizardContext = React.createContext<() => ContextProps>(undefined);

const useWizardContext = () => {
  return React.useContext<() => ContextProps>(WizardContext)();
};

export { WizardContext, useWizardContext };
