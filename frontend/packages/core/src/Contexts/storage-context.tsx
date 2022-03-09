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

export interface StorageState {
  shortLinked: boolean;
  store: HydratedData;
  tempStore: HydratedData;
}

type StoreDataFn = (componentName: string, key: string, data: any, local?: boolean) => void;
type RemoveDataFn = (componentName: string, key: string, local?: boolean) => void;
type RetrieveDataFn = (componentName: string, key: string, defaultData: any) => any;
type ClearDataFn = () => void;

export interface StorageContextProps {
  shortLinked: boolean;
  store?: HydratedData;
  tempStore?: HydratedData;
  functions: {
    storeData: StoreDataFn;
    removeData: RemoveDataFn;
    retrieveData: RetrieveDataFn;
    clearData: ClearDataFn;
  };
}

const StorageContext = React.createContext<StorageContextProps>(undefined);

const useStorageContext = () => {
  return React.useContext<StorageContextProps>(StorageContext);
};

export { StorageContext, useStorageContext };
