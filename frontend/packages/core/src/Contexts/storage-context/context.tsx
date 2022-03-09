import React from "react";

import type { StorageContextProps } from "./types";

const StorageContext = React.createContext<StorageContextProps>(undefined);

const useStorageContext = () => {
  return React.useContext<StorageContextProps>(StorageContext);
};

export { StorageContext, useStorageContext };
