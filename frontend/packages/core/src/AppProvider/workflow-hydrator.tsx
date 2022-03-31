import React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";
import { Grid } from "@material-ui/core";

import { useShortLinkContext } from "../Contexts";
import { WorkflowStorageContext } from "../Contexts/workflow-storage-context";
import { retrieveData } from "../Contexts/workflow-storage-context/helpers";
import workflowStorageContextReducer from "../Contexts/workflow-storage-context/reducer";
import type { WorkflowStorageContextProps } from "../Contexts/workflow-storage-context/types";
import { defaultWorkflowStorageState } from "../Contexts/workflow-storage-context/types";
import { Alert } from "../Feedback";
import styled from "../styled";

interface WorkflowHydratorProps {
  children: React.ReactElement;
  hydrateData: () => IClutch.shortlink.v1.IShareableState[];
  onClear: () => void;
}

const StyledAlert = styled(Alert)({
  zIndex: 1200,
  position: "absolute",
  padding: "6px 8px",
});

/**
 * Hydrator which is a wrapper for workflows
 * Will check on load if there exists any hydrated data for the current workflow
 * If there is it will populate the state and provide an alert above the workflow
 */
const WorkflowHydrator = ({
  hydrateData,
  onClear,
  children,
}: WorkflowHydratorProps): React.ReactElement => {
  const { storeData } = useShortLinkContext();
  const [state, dispatch] = React.useReducer(
    workflowStorageContextReducer,
    defaultWorkflowStorageState
  );

  React.useEffect(() => {
    const hydratedData = hydrateData();
    if (hydratedData) {
      dispatch({ type: "HYDRATE", payload: { data: hydratedData } });
      onClear();
    }
  }, []);

  React.useEffect(() => {
    if (state.tempStore) {
      storeData(state.tempStore);
    }
  }, [state]);

  const storageProviderProps: WorkflowStorageContextProps = {
    shortLinked: state.shortLinked,
    storeData: (componentName: string, key: string, data: unknown, localStorage?: boolean) =>
      dispatch({ type: "STORE_DATA", payload: { componentName, key, data, localStorage } }),
    removeData: (componentName: string, key: string, localStorage?: boolean) =>
      dispatch({ type: "REMOVE_DATA", payload: { componentName, key, localStorage } }),
    retrieveData: (componentName: string, key: string, defaultData?: unknown) =>
      retrieveData(state.store, componentName, key, defaultData),
  };

  return (
    <WorkflowStorageContext.Provider value={storageProviderProps}>
      {state.shortLinked && (
        <Grid container direction="column" alignItems="flex-end">
          <StyledAlert severity="warning">
            Local Workflow Data will not be saved until reload
          </StyledAlert>
        </Grid>
      )}
      {children}
    </WorkflowStorageContext.Provider>
  );
};

export default WorkflowHydrator;
