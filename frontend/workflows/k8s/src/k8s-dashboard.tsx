import React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";
import { client, ClutchError, Error, Paper, styled, Tab, Tabs, useTheme } from "@clutch-sh/core";
import { DataLayoutContext, useDataLayoutManager } from "@clutch-sh/data-layout";
import AppsIcon from "@mui/icons-material/Apps";
import CropFreeIcon from "@mui/icons-material/CropFree";
import DnsOutlinedIcon from "@mui/icons-material/DnsOutlined";
import LoopOutlinedIcon from "@mui/icons-material/LoopOutlined";
import { alpha, Theme } from "@mui/material";

import type { WorkflowProps } from ".";
import CronTable from "./crons-table";
import DeploymentTable from "./deployments-table";
import K8sDashHeader from "./k8s-dash-header";
import K8sDashSearch from "./k8s-dash-search";
import PodTable from "./pods-table";
import StatefulSetTable from "./stateful-sets-table";

const Container = styled("div")(({ theme }: { theme: Theme }) => ({
  flex: 1,
  padding: theme.clutch.layout.gutter,
  display: "flex",
  flexDirection: "column",
  "> *:first-child": {
    margin: theme.spacing("none"),
  },
}));

const Content = styled("div")(({ theme }: { theme: Theme }) => ({
  margin: theme.spacing("lg", "none"),
}));

const PlaceholderTitle = styled("div")(({ theme }: { theme: Theme }) => ({
  paddingBottom: theme.spacing("base"),
  fontWeight: 700,
  fontSize: "22px",
  lineHeight: "28px",
  color: theme.palette.secondary[900],
}));

const PlaceholderText = styled("div")(({ theme }: { theme: Theme }) => ({
  fontWeight: 400,
  fontSize: "16px",
  lineHeight: "22px",
  color: alpha(theme.palette.secondary[900], 0.6),
}));

const PlaceholderWrapper = styled("div")(({ theme }: { theme: Theme }) => ({
  margin: theme.spacing("lg"),
  textAlign: "center",
}));

const Placeholder = () => (
  <Paper>
    <PlaceholderWrapper>
      <PlaceholderTitle>There is nothing to display here</PlaceholderTitle>
      <PlaceholderText>Please enter a namespace and clientset to proceed.</PlaceholderText>
    </PlaceholderWrapper>
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
  const theme = useTheme();
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
        {!theme.clutch.useWorkflowLayout && <K8sDashHeader />}
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
