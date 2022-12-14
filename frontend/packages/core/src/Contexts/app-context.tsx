import React from "react";

import type { Workflow } from "../AppProvider/workflow";

export type HeaderItems = "NPS";

export interface TriggeredHeaderData {
  [key: string]: {
    open: boolean;
    [key: string]: unknown;
  };
}

interface ContextProps {
  workflows: Workflow[];
  triggerHeaderItem?: (item: HeaderItems, open: boolean, data?: unknown) => void;
  triggeredHeaderData?: TriggeredHeaderData;
}

const ApplicationContext = React.createContext<ContextProps>(undefined);

const useAppContext = () => {
  return React.useContext<ContextProps>(ApplicationContext);
};

export { ApplicationContext, useAppContext };
