import React from "react";

import { Alert } from "../Feedback";
import Grid from "../grid";
import type { AppConfiguration } from "../Types";

const LayoutWithNotifications = ({
  children,
  config,
  workflow,
}: {
  children: React.ReactNode;
  config: AppConfiguration;
  workflow?: string;
}) => {
  const {
    banners: { perWorkflow, multiWorkflow },
  } = config;

  const showAlertPerWorkflow =
    workflow && perWorkflow[workflow] && !perWorkflow[workflow]?.dismissed;
  const showAlertMultiWorkflow =
    workflow && multiWorkflow?.workflows.includes(workflow) && !multiWorkflow.dismissed;

  return (
    <>
      {config && config.banners && (
        <Grid container justifyContent="center" pt={2} pb={1} px={3}>
          <Grid item xs>
            {showAlertMultiWorkflow && (
              <Alert severity={multiWorkflow?.severity} title={multiWorkflow?.title} elevation={6}>
                {multiWorkflow?.message}
              </Alert>
            )}
            {showAlertPerWorkflow && (
              <Alert
                severity={perWorkflow[workflow]?.severity}
                title={perWorkflow[workflow]?.title}
                elevation={6}
              >
                {perWorkflow[workflow]?.message}
              </Alert>
            )}
          </Grid>
        </Grid>
      )}
      {children}
    </>
  );
};

export default LayoutWithNotifications;
