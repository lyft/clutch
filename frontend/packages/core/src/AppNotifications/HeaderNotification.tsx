import React from "react";
import styled from "@emotion/styled";
import isEmpty from "lodash/isEmpty";

import { Alert } from "../Feedback";
import Grid from "../grid";
import { Link as LinkComponent } from "../link";
import type { AppBanners } from "../Types";

const StyledAlert = styled(Alert)({
  padding: "8px 16px 8px 16px",
  justifyContent: "center",
  alignItems: "center",
});

const StyledAlertContent = styled.div({
  display: "flex",
  maxHeight: "40px",
  overflowY: "auto",
});

const StyledMessage = styled.div({
  flexWrap: "wrap",
});

const StyledLink = styled.div({
  marginLeft: "10px",
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
        <Grid item xs={4}>
          <StyledAlert
            severity={headerBannerData?.severity || "info"}
            elevation={6}
            onClose={onDismissAlertHeader}
          >
            <StyledAlertContent>
              <StyledMessage>{headerBannerData.message}</StyledMessage>
              {headerBannerData?.link && headerBannerData?.linkText && (
                <StyledLink>
                  <LinkComponent href={headerBannerData?.link}>
                    {headerBannerData?.linkText}
                  </LinkComponent>
                </StyledLink>
              )}
            </StyledAlertContent>
          </StyledAlert>
        </Grid>
      )}
    </>
  );
};

export default HeaderNotification;
