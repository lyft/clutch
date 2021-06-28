import React from "react";

export interface ContextProps {
  displayWarnings: (warnings: string[]) => void;
  onBack: (params?: { toOrigin?: boolean }) => void;
  onSubmit: () => void;
  setOnSubmit: (f: (...args: any[]) => void) => void;
  setIsLoading: (isLoading: boolean) => void;
  setHasError: (hasError: boolean) => void;
}

const WizardContext = React.createContext<() => ContextProps>(undefined);

const useWizardContext = () => {
  return React.useContext<() => ContextProps>(WizardContext)();
};

export { WizardContext, useWizardContext };
