import React from "react";

import type { HydratedData } from "./workflow-storage-context/types";

export interface ShortLinkContextProps {
  removeWorkflowSession: () => void;
  retrieveWorkflowSession: () => HydratedData;
  storeWorkflowSession: (data: HydratedData) => void;
}

const ShortLinkContext = React.createContext<ShortLinkContextProps>(undefined);

const useShortLinkContext = () => {
  return React.useContext<ShortLinkContextProps>(ShortLinkContext);
};

export { ShortLinkContext, useShortLinkContext };
