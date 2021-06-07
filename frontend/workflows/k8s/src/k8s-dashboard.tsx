import React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";
import { client, ClutchError, Error, Paper, Tab, Tabs } from "@clutch-sh/core";
import { DataLayoutContext, useDataLayoutManager } from "@clutch-sh/data-layout";
import styled from "@emotion/styled";
import AppsIcon from "@material-ui/icons/Apps";
import CropFreeIcon from "@material-ui/icons/CropFree";
import DnsOutlinedIcon from "@material-ui/icons/DnsOutlined";
import LoopOutlinedIcon from "@material-ui/icons/LoopOutlined";

import type { WorkflowProps } from ".";
import CronTable from "./crons-table";
import DeploymentTable from "./deployments-table";
import K8sDashHeader from "./k8s-dash-header";
import K8sDashSearch from "./k8s-dash-search";
import PodTable from "./pods-table";
import StatefulSetTable from "./stateful-sets-table";

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
    deploymentListData: {
      deps: ["inputData"],
      hydrator: inputData => {
        return client
          .post("/v1/k8s/listDeployments", {
            ...defaultRequestData(inputData),
            options: {
              labels: {},
            },
          } as IClutch.k8s.v1.IListDeploymentsRequest)
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
    cronListData: {
      deps: ["inputData"],
      hydrator: inputData => {
        return client
          .post("/v1/k8s/listCronJobs", {
            ...defaultRequestData(inputData),
            options: {
              labels: {},
            },
          } as IClutch.k8s.v1.IListCronJobsRequest)
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
    statefulSetListData: {
      deps: ["inputData"],
      hydrator: inputData => {
        return client
          .post("/v1/k8s/listStatefulSets", {
            ...defaultRequestData(inputData),
            options: {
              labels: {},
            },
          } as IClutch.k8s.v1.IListStatefulSetsRequest)
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
    dataLayoutManager.hydrate("deploymentListData");
    dataLayoutManager.hydrate("cronListData");
    dataLayoutManager.hydrate("statefulSetListData");
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
                <Tab startAdornment={<CropFreeIcon />} label="Deployments">
                  <DeploymentTable />
                </Tab>
                <Tab startAdornment={<LoopOutlinedIcon />} label="Cron Jobs">
                  <CronTable />
                </Tab>
                <Tab startAdornment={<DnsOutlinedIcon />} label="Stateful Sets">
                  <StatefulSetTable />
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
