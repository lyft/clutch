import React from "react";

export interface NavigationProps {
  toOrigin?: boolean;
  keepSearch?: boolean;
}

export interface ContextProps {
  displayWarnings: (warnings: string[]) => void;
  onBack: (params?: NavigationProps) => void;
  onNext: (params?: NavigationProps) => void;
  onSubmit: () => void;
  setOnSubmit: (f: (...args: any[]) => void) => void;
  setIsLoading: (isLoading: boolean) => void;
  setHasError: (hasError: boolean) => void;
  onComplete: (id: string) => void;
}

const WizardContext = React.createContext<() => ContextProps>(undefined);

const useWizardContext = () => {
  return React.useContext<() => ContextProps>(WizardContext)();
};

export { WizardContext, useWizardContext };
