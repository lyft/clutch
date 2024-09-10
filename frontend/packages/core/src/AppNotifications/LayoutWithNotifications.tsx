import React from "react";
import { useLocation } from "react-router-dom";
import isEmpty from "lodash/isEmpty";

import { Alert } from "../Feedback";
import Grid from "../grid";
import { Link as LinkComponent } from "../link";
import type { AppBanners } from "../Types";
import { checkPathMatchList } from "../utils";

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

  const location = useLocation();

  checkPathMatchList(location?.pathname, perWorkflowData[workflow]?.paths);

  const hasPerWorkflowAlert =
    workflow && perWorkflowData[workflow] && !perWorkflowData[workflow]?.dismissed;
  const showAlertPerWorkflow = !isEmpty(perWorkflowData[workflow]?.paths)
    ? hasPerWorkflowAlert && perWorkflowData[workflow]?.paths?.includes(location.pathname)
    : hasPerWorkflowAlert;

  const showAlertMultiWorkflow =
    showAlertPerWorkflow || perWorkflowData[workflow]?.dismissed
      ? false
      : workflow &&
        multiWorkflowData?.workflows?.includes(workflow) &&
        !multiWorkflowData?.dismissed;

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
                severity={perWorkflowData[workflow]?.severity || "info"}
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
                severity={multiWorkflowData?.severity || "info"}
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
