import React from "react";
import isEmpty from "lodash/isEmpty";

import { Alert } from "../Feedback";
import Grid from "../grid";
import { Link as LinkComponent } from "../link";
import type { AppBanners } from "../Types";

interface LayoutWithNotificationsProps {
  bannersData: AppBanners;
  onDismissAlert: (updatedData: AppBanners) => void;
  children: React.ReactNode;
  workflow?: string;
}

const LayoutWithNotifications = ({
  bannersData,
  onDismissAlert,
  children,
  workflow,
}: LayoutWithNotificationsProps) => {
  const perWorkflowData = bannersData?.perWorkflow;
  const multiWorkflowData = bannersData?.multiWorkflow;

  const showAlertPerWorkflow =
    workflow && perWorkflowData[workflow] && !perWorkflowData[workflow]?.dismissed;
  const showAlertMultiWorkflow =
    workflow && multiWorkflowData?.workflows?.includes(workflow) && !multiWorkflowData?.dismissed;

  const onDismissAlertPerWorkflow = () => {
    onDismissAlert({
      ...bannersData,
      perWorkflow: {
        ...perWorkflowData,
        [workflow]: { ...perWorkflowData[workflow], dismissed: true },
      },
    });
  };

  const onDismissAlertMultiWorkflow = () => {
    onDismissAlert({
      ...bannersData,
      multiWorkflow: {
        ...multiWorkflowData,
        dismissed: true,
      },
    });
  };

  const showContainer = !isEmpty(perWorkflowData) || !isEmpty(multiWorkflowData);

  return (
    <>
      {showContainer && (
        <Grid container justifyContent="center" pt={2} pb={1} px={3}>
          <Grid item xs>
            {showAlertPerWorkflow && (
              <Alert
                severity={perWorkflowData[workflow]?.severity}
                title={perWorkflowData[workflow]?.title}
                elevation={6}
                onClose={onDismissAlertPerWorkflow}
              >
                {perWorkflowData[workflow]?.message}
                {perWorkflowData[workflow]?.link && perWorkflowData[workflow]?.linkText && (
                  <LinkComponent href={perWorkflowData[workflow]?.link}>
                    {perWorkflowData[workflow]?.linkText}
                  </LinkComponent>
                )}
              </Alert>
            )}
            {showAlertMultiWorkflow && !showAlertPerWorkflow && (
              <Alert
                severity={multiWorkflowData?.severity}
                title={multiWorkflowData?.title}
                elevation={6}
                onClose={onDismissAlertMultiWorkflow}
              >
                {multiWorkflowData?.message}
                {multiWorkflowData?.link && multiWorkflowData?.linkText && (
                  <LinkComponent href={multiWorkflowData?.link}>
                    {multiWorkflowData?.linkText}
                  </LinkComponent>
                )}
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
