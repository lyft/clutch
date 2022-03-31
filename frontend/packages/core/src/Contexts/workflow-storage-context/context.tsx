import React from "react";

import type { WorkflowStorageContextProps } from "./types";

const WorkflowStorageContext = React.createContext<WorkflowStorageContextProps>(undefined);

const useWorkflowStorageContext = () => {
  return React.useContext<WorkflowStorageContextProps>(WorkflowStorageContext);
};

export { WorkflowStorageContext, useWorkflowStorageContext };
