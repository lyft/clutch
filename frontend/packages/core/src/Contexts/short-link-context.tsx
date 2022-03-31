import React from "react";

import type { HydratedData } from "./workflow-storage-context/types";

export interface ShortLinkContextProps {
  removeData: () => void;
  retrieveData: () => HydratedData;
  storeData: (data: HydratedData) => void;
}

const ShortLinkContext = React.createContext<ShortLinkContextProps>(undefined);

const useShortLinkContext = () => {
  return React.useContext<ShortLinkContextProps>(ShortLinkContext);
};

export { ShortLinkContext, useShortLinkContext };
