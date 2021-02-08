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
  assign: (value: object) => void;
  updateData: (dataKey: string, value: unknown) => void;
  hydrate: () => void;
  value: any;
  displayValue: () => any;
  isLoading: boolean;
  error: string;
  setError: (error: string) => void;
  setLoading: () => void;
}

const useDataLayout = (key: string, opts?: object): DataLayout => {
  const manager = useManagerContext();

  const options = { hydrate: true, ...opts };

  if (!Object.keys(manager.state).includes(key)) {
    throw new Error(`Non-existant data layout key: ${key}`);
  }

  React.useEffect(() => {
    manager.update(key, { error: "", isLoading: false });
  }, []);

  React.useEffect(() => {
    if (options.hydrate && (!manager.state[key].cache || _.isEmpty(manager.state[key].data))) {
      manager.hydrate(key);
    }
  }, [key]);

  return {
    assign: value => value !== manager.state[key].data && manager.assign(key, value),
    updateData: (dataKey, value) => updateData(manager, key, dataKey, value),
    hydrate: () => manager.hydrate(key),
    value: manager.state[key].data,
    displayValue: () =>
      manager.state[key].data?.toJSON ? manager.state[key].data?.toJSON() : manager.state[key].data,
    isLoading: manager.state[key].isLoading,
    error: manager.state[key].error,
    setError: error =>
      error !== manager.state[key].error && manager.update(key, { error, isLoading: false }),
    setLoading: () => manager.update(key, { isLoading: true }),
  };
};

export default useDataLayout;
