import React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";
import { client, ClutchError, Error, Paper, Tab, Tabs } from "@clutch-sh/core";
import { DataLayoutContext, useDataLayoutManager } from "@clutch-sh/data-layout";
import styled from "@emotion/styled";
import AppsIcon from "@material-ui/icons/Apps";

import type { WorkflowProps } from ".";
import K8sDashHeader from "./k8s-dash-header";
import K8sDashSearch from "./k8s-dash-search";
import PodTable from "./pods-table";

const Container = styled.div({
  flex: 1,
  margin: "32px",
  display: "flex",
  flexDirection: "column",
  "> *:first-child": {
    margin: "0",
  },
});

const Content = styled.div({
  margin: "32px 0",
});

const PlaceholderTitle = styled.div({
  paddingBottom: "16px",
  fontWeight: 700,
  fontSize: "22px",
  lineHeight: "28px",
  color: "0D1030",
});

const PlaceholderText = styled.div({
  fontWeight: 400,
  fontSize: "16px",
  lineHeight: "22px",
  color: "rgba(13, 16, 48, 0.6)",
});

const Placeholder = () => (
  <Paper>
    <div style={{ margin: "32px", textAlign: "center" }}>
      <PlaceholderTitle>There is nothing to display here</PlaceholderTitle>
      <PlaceholderText>Please enter a namespace and clientset to proceed.</PlaceholderText>
    </div>
  </Paper>
);

const defaultRequestData = inputData => {
  return {
    clientset: inputData.clientset,
    cluster: inputData.clientset,
    namespace: inputData.namespace,
  };
};

const KubeDashboard: React.FC<WorkflowProps> = () => {
  const [error, setError] = React.useState<ClutchError | undefined>(undefined);

  const dataLayout = {
    inputData: {},
    podListData: {
      deps: ["inputData"],
      hydrator: inputData => {
        return client
          .post("/v1/k8s/listPods", {
            ...defaultRequestData(inputData),
            options: {
              labels: {},
            },
          } as IClutch.k8s.v1.IListPodsRequest)
          .then(response => {
            return response?.data;
          })
          .catch((err: ClutchError) => {
            setError(existingError => {
              if (existingError === undefined) {
                return err;
              }
              return existingError;
            });
          });
      },
    },
  };
  const dataLayoutManager = useDataLayoutManager(dataLayout);

  const handleSubmit = (namespace, clientset) => {
    dataLayoutManager.assign("inputData", { namespace, clientset });
    dataLayoutManager.hydrate("podListData");
    setError(undefined);
  };

  const state = dataLayoutManager.state as any;
  const hasData = state.inputData?.data?.namespace;

  return (
    <DataLayoutContext.Provider value={dataLayoutManager}>
      <Container>
        <K8sDashHeader />
        <K8sDashSearch onSubmit={(namespace, clientset) => handleSubmit(namespace, clientset)} />
        <Content>
          {error !== undefined ? (
            <Error subject={error} />
          ) : !hasData ? (
            <Placeholder />
          ) : (
            <Paper>
              <Tabs variant="fullWidth">
                <Tab startAdornment={<AppsIcon />} label="Pods">
                  <PodTable />
                </Tab>
              </Tabs>
            </Paper>
          )}
        </Content>
      </Container>
    </DataLayoutContext.Provider>
  );
};

export default KubeDashboard;
