import React from "react";

import type { Workflow } from "../AppProvider/workflow";

interface ContextProps {
  workflows: Workflow[];
  appConfig: {[key: string]: any};
}
const ApplicationContext = React.createContext<ContextProps>(undefined);

const useAppContext = () => {
  return React.useContext<ContextProps>(ApplicationContext);
};

export { ApplicationContext, useAppContext };
