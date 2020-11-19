import React from "react";
import styled from "@emotion/styled";
import { Icon } from "@material-ui/core";
import NotificationsIcon from "@material-ui/icons/Notifications";

const NotifcationsStyled = styled(Icon)`
  color: #ffffff;
  opacity: 0.87;
`;

const Notifications: React.FC = () => {
  return (
    <NotifcationsStyled>
      <NotificationsIcon />
    </NotifcationsStyled>
  );
};

export default Notifications;