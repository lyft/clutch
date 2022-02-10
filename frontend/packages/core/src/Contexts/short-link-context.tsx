import React from "react";

export interface HydrateData {
  data: any;
}

interface ContextProps {
  hydration?: HydrateData;
}

const ShortLinkContext = React.createContext<ContextProps>(undefined);

const useShortLinkContext = () => {
  return React.useContext<ContextProps>(ShortLinkContext);
};

export { ShortLinkContext, useShortLinkContext };
