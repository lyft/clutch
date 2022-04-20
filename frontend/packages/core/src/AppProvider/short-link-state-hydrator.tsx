import React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";
import { Grid } from "@material-ui/core";

import { WorkflowStorageContext } from "../Contexts/workflow-storage-context";
import { retrieveData } from "../Contexts/workflow-storage-context/helpers";
import workflowStorageContextReducer from "../Contexts/workflow-storage-context/reducer";
import type { WorkflowStorageContextProps } from "../Contexts/workflow-storage-context/types";
import { defaultWorkflowStorageState } from "../Contexts/workflow-storage-context/types";
import { Alert } from "../Feedback";
import styled from "../styled";

interface ShortLinkStateHydratorProps {
  children: React.ReactElement;
  /** Data from ShortLink API to be hydrated into the  */
  sharedState: IClutch.shortlink.v1.IShareableState[];
  /** Used to clear temporary storage variable in the AppProvider */
  clearTemporaryState: () => void;
}

/** Allows the alert to float on top of all other components on the page */
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
const ShortLinkStateHydrator = ({
  children,
  clearTemporaryState,
  sharedState,
}: ShortLinkStateHydratorProps): React.ReactElement => {
  const [state, dispatch] = React.useReducer(
    workflowStorageContextReducer,
    defaultWorkflowStorageState
  );

  React.useEffect(() => {
    if (sharedState && sharedState.length) {
      dispatch({ type: "HYDRATE", payload: { data: sharedState } });
      clearTemporaryState();
    }
  }, [sharedState]);

  const storageProviderProps: WorkflowStorageContextProps = {
    /**
     * Boolean representing whether the state has been hydrated
     */
    fromShortLink: state.fromShortLink,
    /**
     * StoreData context function which will allow a component to write data for use in shortlinks as well as
     * store locally
     * @param componentName Name of the component that data is being stored under
     * @param key A lookup key used to reference the specific data set being stored
     * @param data The data being stored
     * @param localStorage Optional boolean on whether to write data to the localStorage as well
     * @returns void
     */
    storeData: (componentName: string, key: string, data: unknown, localStorage?: boolean) =>
      dispatch({ type: "STORE_DATA", payload: { componentName, key, data, localStorage } }),
    /**
     * RemoveData context function which will allow a component to remove data from use in shortlinks as well
     * locally if preferred
     * @param componentName Name of the component that data is being removed under
     * @param key A lookup key used to reference the specific data set being removed
     * @param localStorage Optional boolean on whether to remove data from localStorage as well
     * @returns
     */
    removeData: (componentName: string, key: string, localStorage?: boolean) =>
      dispatch({ type: "REMOVE_DATA", payload: { componentName, key, localStorage } }),
    /**
     * RetrieveData context function which will allow a component to retrieve data from a hydrated short link and
     * barring that will return data from local storage if available
     * @param componentName Name of the component that data is being retrieved under
     * @param key A lookup key used to reference the specific data set being retrieved
     * @param defaultData Optional set of data returned if nothing is found in the hydrator or localStorage
     * @returns
     */
    retrieveData: (componentName: string, key: string, defaultData?: unknown) =>
      retrieveData(state.workflowStore, componentName, key, defaultData),
  };

  return (
    <WorkflowStorageContext.Provider value={storageProviderProps}>
      {state.fromShortLink && (
        <Grid container direction="column" alignItems="flex-end">
          <StyledAlert severity="warning">
            Loaded shared state. Any local changes will not be preserved.
          </StyledAlert>
        </Grid>
      )}
      {children}
    </WorkflowStorageContext.Provider>
  );
};

export default ShortLinkStateHydrator;
