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
  hydrateStore?: HydratedData;
  tempHydrateStore?: HydratedData;
  data: {
    store: (componentName: string, key: string, data: any, local?: boolean) => void;
    localStore: (key: string, data: any) => void;
    retrieve: (componentName: string, key: string, defaultData?: any) => any;
    remove: (componentName: string, key: string, local?: boolean) => void;
  };
}

const StorageContext = React.createContext<ContextProps>(undefined);

const useStorageContext = () => {
  return React.useContext<ContextProps>(StorageContext);
};

export { StorageContext, useStorageContext };
