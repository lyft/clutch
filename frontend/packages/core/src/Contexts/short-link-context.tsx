import React from "react";
import type { GlobalProjectState } from "@clutch-sh/project-selector";

export interface HydrateData {
  route: string;
  data: HydratedData;
}

export interface HydratedData {
  dash?: {
    state?: GlobalProjectState;
    splitEvents?: boolean;
  };
}

interface ContextProps {
  hydration?: HydratedData;
}

const ShortLinkContext = React.createContext<ContextProps>(undefined);

const useShortLinkContext = () => {
  return React.useContext<ContextProps>(ShortLinkContext);
};

export { ShortLinkContext, useShortLinkContext };
