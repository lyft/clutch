import React from "react";
import _ from "lodash";

import { useManagerContext } from "./context";
import type { DataManager } from "./manager";

const updateData = (manager: DataManager, layoutKey: string, key: string, value: any) => {
  const { data } = manager.state[layoutKey];
  _.set(data, key, value);
  manager.update(layoutKey, { data });
};

interface DataLayout {
  /**
   * Overwrite existing data for the layout with the specified data.
   */
  assign: (data: object) => void;
  /**
   * Update existing data with the specified key, value pair.
   */
  updateData: (key: string, value: object) => void;
  /**
   * Hydrate the layout data with the return value of the specified function.
   */
  hydrate: () => void;
  /**
   * The raw data of the layout.
   */
  value: object;
  /**
   * Returns a json representation of the layout's data if possible, otherwise returns the raw data.
   */
  displayValue: () => object;
  /**
   * Loading state of the layout. This is true when data is being hydrated.
   */
  isLoading: boolean;
  /**
   * Error state of the layout. This will be a message containing the error encountered when trying to hydrate the lyaout.
   */
  error: string;
}

/**
 * Use a registered data layout.
 * 
 * If a hydrate function has been specified this and the layout's data has not been set this will populate it's data
 * on the first invocation. If the layout has a cache key set to true and also has existing data hydrate will not be invoked.

 * @param key The name of the layout registered with the manager.
 */
const useDataLayout = (key: string): DataLayout => {
  const manager = useManagerContext();

  const options = { hydrate: true };

  if (!Object.keys(manager.state).includes(key)) {
    throw new Error(`Non-existant data layout key: ${key}`);
  }

  // n.b. reset error and loading state on load.
  // This prevents previous errors from rendering until hydration is finished.
  React.useEffect(() => {
    manager.update(key, { error: "", isLoading: false });
  }, []);

  React.useEffect(() => {
    if (
      options.hydrate !== undefined &&
      !(manager.state[key].cache && !_.isEmpty(manager.state[key].data))
    ) {
      manager.hydrate(key);
    }
  }, [key]);

  return {
    assign: data => data !== manager.state[key].data && manager.assign(key, data),
    updateData: (dataKey, value) => updateData(manager, key, dataKey, value),
    hydrate: () => manager.hydrate(key),
    value: manager.state[key].data,
    displayValue: () =>
      manager.state[key].data?.toJSON ? manager.state[key].data.toJSON() : manager.state[key].data,
    isLoading: manager.state[key].isLoading,
    error: manager.state[key].error,
  };
};

export default useDataLayout;
