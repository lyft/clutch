import React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";

import { useShortLinkContext } from "../Contexts/shortlink-context";
import { WorkflowStorageContext } from "../Contexts/workflow-storage-context";
import { retrieveData } from "../Contexts/workflow-storage-context/helpers";
import workflowStorageContextReducer from "../Contexts/workflow-storage-context/reducer";
import type { WorkflowStorageContextProps } from "../Contexts/workflow-storage-context/types";
import { defaultWorkflowStorageState } from "../Contexts/workflow-storage-context/types";
import { Toast } from "../Feedback";

interface ShortLinkHydratorProps {
  hydrate: () => IClutch.shortlink.v1.IShareableState[] | null;
  onClear: () => void;
  children: React.ReactElement;
}

/**
 * Hydrator which is a wrapper for workflows
 * Will check on load if there exists any hydrated data for the current workflow
 * If there is it will populate the state and provide a toast
 */
const ShortLinkHydrator = ({
  hydrate,
  onClear,
  children,
}: ShortLinkHydratorProps): React.ReactElement => {
  const { storeData } = useShortLinkContext();
  const [workflowStorageState, dispatch] = React.useReducer(
    workflowStorageContextReducer,
    defaultWorkflowStorageState
  );

  React.useEffect(() => {
    const data = hydrate();

    if (data) {
      dispatch({ type: "HYDRATE", payload: { data } });
      onClear();
    }
  }, []);

  React.useEffect(() => {
    if (workflowStorageState.tempStore) {
      storeData(workflowStorageState.tempStore);
    }
  }, [workflowStorageState]);

  const workflowStorageProviderProps: WorkflowStorageContextProps = {
    shortLinked: workflowStorageState.shortLinked,
    storeData: (componentName: string, key: string, data: any, localStorage?: boolean) =>
      dispatch({ type: "STORE_DATA", payload: { componentName, key, data, localStorage } }),
    removeData: (componentName: string, key: string, localStorage?: boolean) =>
      dispatch({ type: "REMOVE_DATA", payload: { componentName, key, localStorage } }),
    retrieveData: (componentName: string, key: string, defaultData?: any) =>
      retrieveData(workflowStorageState.store, componentName, key, defaultData),
  };

  return (
    <WorkflowStorageContext.Provider value={workflowStorageProviderProps}>
      {workflowStorageState.shortLinked && (
        <Toast title="Visited Short Link">Local Workflow Data will not be saved until reload</Toast>
      )}
      {children}
    </WorkflowStorageContext.Provider>
  );
};

export default ShortLinkHydrator;
