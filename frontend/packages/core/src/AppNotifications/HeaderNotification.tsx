import React from "react";
import styled from "@emotion/styled";
import isEmpty from "lodash/isEmpty";

import { Alert } from "../Feedback";
import Grid from "../grid";
import { Link as LinkComponent } from "../link";
import type { AppBanners } from "../Types";

const StyledAlert = styled(Alert)({
  padding: "8px 16px 8px 16px",
});

interface HeaderNotificationProps {
  bannersData: AppBanners;
  onDismissAlert: (updatedData: AppBanners) => void;
}

const HeaderNotification = ({ bannersData, onDismissAlert }: HeaderNotificationProps) => {
  const headerBannerData = bannersData?.header;

  const onDismissAlertHeader = () => {
    onDismissAlert({
      ...bannersData,
      header: {
        ...headerBannerData,
        dismissed: true,
      },
    });
  };

  return (
    <>
      {!isEmpty(headerBannerData) && !headerBannerData?.dismissed && (
        <Grid item xs={3}>
          <StyledAlert
            severity={headerBannerData?.severity}
            title={headerBannerData?.title}
            elevation={6}
            onClose={onDismissAlertHeader}
          >
            {headerBannerData?.message}
            {headerBannerData?.link && <LinkComponent>{headerBannerData?.link}</LinkComponent>}
          </StyledAlert>
        </Grid>
      )}
    </>
  );
};

export default HeaderNotification;
