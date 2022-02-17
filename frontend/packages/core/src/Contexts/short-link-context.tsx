import React from "react";

export interface HydrateData {
  route: string;
  data: HydratedData;
}

export interface HydratedData {
  [key: string]: {
    [key: string]: any;
  };
}

interface ContextProps {
  hydration?: HydratedData;
  tempHydrateStore?: HydrateData;
  storeHydration?: (data: any) => void;
}

const ShortLinkContext = React.createContext<ContextProps>(undefined);

const useShortLinkContext = () => {
  return React.useContext<ContextProps>(ShortLinkContext);
};

export { ShortLinkContext, useShortLinkContext };
