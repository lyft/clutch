import React from "react";

export interface ContextProps {
  displayWarnings: (warnings: string[]) => void;
  onBack: () => void;
  onSubmit: () => void;
  setOnSubmit: (f: (...args: any[]) => void) => void;
  setIsLoading: (isLoading: boolean) => void;
}

const WizardContext = React.createContext<() => ContextProps>(undefined);

const useWizardContext = () => {
  return React.useContext<() => ContextProps>(WizardContext)();
};

export { WizardContext, useWizardContext };
