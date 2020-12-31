import * as React from "react";
import styled from "@emotion/styled";
import { Grid } from "@material-ui/core";
import MuiSuccessIcon from "@material-ui/icons/CheckCircle";
import MuiErrorIcon from "@material-ui/icons/Error";
import MuiInfoIcon from "@material-ui/icons/Info";
import MuiWarningIcon from "@material-ui/icons/Warning";
import type { AlertProps as MuiAlertProps } from "@material-ui/lab";
import { Alert as MuiAlert, AlertTitle as MuiAlertTitle } from "@material-ui/lab";

const backgroundColors = {
  error: "linear-gradient(to right, #DB3615 8px, #FDE9E7 0%)",
  info: "linear-gradient(to right, #3548D4 8px, #EBEDFB 0%)",
  success: "linear-gradient(to right, #1E942E 8px, #E6F7EB 0%)",
  warning: "linear-gradient(to right, #FFCC80 8px, #FFFDE6 0%)",
};

const StyledAlert = styled(MuiAlert)(
  {
    borderRadius: "8px",
    padding: "16px",
    paddingBottom: "20px",
    color: "rgba(13, 16, 48, 0.6)",
    fontSize: "14px",
    ".MuiAlert-icon": {
      marginRight: "16px",
      padding: "0",
    },
    ".MuiAlert-message": {
      padding: "0",
      ".MuiAlertTitle-root": {
        marginBottom: "0",
        color: "#0D1030",
      },
    },
  },
  props => ({
    background: backgroundColors[props.severity],
  })
);

const ErrorIcon = styled(MuiErrorIcon)({
  color: "#db3716",
});

const InfoIcon = styled(MuiInfoIcon)({
  color: "#3548d4",
});

const SuccessIcon = styled(MuiSuccessIcon)({
  color: "#1e942d",
});

const WarningIcon = styled(MuiWarningIcon)({
  color: "#ffcc80",
});

const AlertTitle = styled(MuiAlertTitle)({
  color: "#0D1030",
  fontWeight: 600,
  fontSize: "16px",
});

const iconMappings = {
  error: <ErrorIcon />,
  info: <InfoIcon />,
  success: <SuccessIcon />,
  warning: <WarningIcon />,
};
export interface AlertProps extends Pick<MuiAlertProps, "severity" | "action"> {
  title?: React.ReactNode;
}

export const Alert: React.FC<AlertProps> = ({ severity = "info", title, children, ...props }) => (
  <StyledAlert severity={severity} iconMapping={iconMappings} {...props}>
    {title && <AlertTitle>{title}</AlertTitle>}
    {children}
  </StyledAlert>
);

export interface AlertPanelProps {
  direction?: "row" | "column";
  children: React.ReactElement<AlertProps> | React.ReactElement<AlertProps>[];
}

export const AlertPanel = ({ direction = "column", children }: AlertPanelProps) => (
  <Grid container direction={direction} justify="center" alignContent="space-between" wrap="nowrap">
    {children}
  </Grid>
);
