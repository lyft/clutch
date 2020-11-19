import React from "react";
import styled from "@emotion/styled";
import { Icon } from "@material-ui/core";
import NotificationsIcon from "@material-ui/icons/Notifications";

const StyledNotificationsIcon = styled(Icon)`
  color: #ffffff;
  opacity: 0.87;
`;

const Notifications: React.FC = () => {
  return (
    <StyledNotificationsIcon>
      <NotificationsIcon />
    </StyledNotificationsIcon>
  );
};

export default Notifications;