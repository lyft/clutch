import React from "react";

import type { Workflow } from "../AppProvider/workflow";

export type HeaderLinks = "NPS";

export interface HeaderItems {
  [key: string]: {
    open: boolean;
    [key: string]: any;
  };
}

interface ContextProps {
  workflows: Workflow[];
  headerLink: (item: HeaderLinks, data?: any) => void;
  headerItems: HeaderItems;
}
const ApplicationContext = React.createContext<ContextProps>(undefined);

const useAppContext = () => {
  return React.useContext<ContextProps>(ApplicationContext);
};

export { ApplicationContext, useAppContext };
