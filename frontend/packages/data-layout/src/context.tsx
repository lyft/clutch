import React from "react";

import type { DataManager } from "./manager";

const DataLayoutContext = React.createContext<DataManager>(undefined);

const useManagerContext = (): DataManager => {
  return React.useContext<DataManager>(DataLayoutContext);
};

export { DataLayoutContext, useManagerContext };
