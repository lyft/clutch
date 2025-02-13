import React from "react";

export interface WizardNavigationProps {
  toOrigin?: boolean;
  keepSearch?: boolean;
}

export interface ConfirmActionProps {
  title: string;
  description: React.ReactNode | string;
  actionLabel?: string;
  confirmationText: string;
  onConfirm: () => void;
}

export interface ContextProps {
  displayWarnings: (warnings: string[]) => void;
  onBack: (params?: WizardNavigationProps) => void;
  onNext?: (params?: WizardNavigationProps) => void;
  onSubmit: () => void;
  setOnSubmit: (f: (...args: any[]) => void) => void;
  setIsLoading: (isLoading: boolean) => void;
  setHasError: (hasError: boolean) => void;
  setIsComplete?: (isComplete: boolean) => void;
  showConfirmAction?: (props: ConfirmActionProps) => void;
}

const WizardContext = React.createContext<() => ContextProps>(undefined);

const useWizardContext = () => {
  return React.useContext<() => ContextProps>(WizardContext)();
};

export { WizardContext, useWizardContext };
