import * as React from "react";
import MuiSuccessIcon from "@mui/icons-material/CheckCircle";
import MuiErrorIcon from "@mui/icons-material/Error";
import MuiInfoIcon from "@mui/icons-material/Info";
import MuiWarningIcon from "@mui/icons-material/Warning";
import type { AlertProps as MuiAlertProps } from "@mui/lab";
import { Alert as MuiAlert, AlertTitle as MuiAlertTitle, Grid, useTheme } from "@mui/material";

import styled from "../styled";

const StyledAlert = styled(MuiAlert)<{ severity: MuiAlertProps["severity"] }>(
  () => {
    const theme = useTheme();
    return {
      borderRadius: "8px",
      padding: "16px",
      paddingLeft: "24px",
      paddingBottom: "20px",
      color: `${theme.palette.secondary[700]}66`,
      fontSize: "14px",
      overflow: "auto",
      ".MuiAlert-icon": {
        marginRight: "16px",
        padding: "0",
      },
      ".MuiAlert-message": {
        maxWidth: "calc(100% - 40px)",
        width: "100%",
        padding: "0",
        ".MuiAlertTitle-root": {
          marginBottom: "0",
          color: theme.palette.secondary[900],
        },
      },
    };
  },
  props => () => {
    const theme = useTheme();
    const backgroundColors = {
      error: `linear-gradient(to right, ${theme.palette.error[600]} 8px, ${theme.palette.error[100]} 0%)`,
      info: `linear-gradient(to right, ${theme.palette.primary[600]} 8px, ${theme.palette.primary[200]} 0%)`,
      success: `linear-gradient(to right, ${theme.palette.success[500]} 8px, ${theme.palette.success[100]} 0%)`,
      warning: `linear-gradient(to right, ${theme.palette.warning[500]} 8px, ${theme.palette.warning[100]} 0%)`,
    };

    return {
      background: backgroundColors[props.severity],
    };
  }
);

const ErrorIcon = styled(MuiErrorIcon)(() => {
  const theme = useTheme();
  return {
    color: theme.palette.error[700],
  };
});

const InfoIcon = styled(MuiInfoIcon)(() => {
  const theme = useTheme();
  return {
    color: theme.palette.primary[600],
  };
});

const SuccessIcon = styled(MuiSuccessIcon)(() => {
  const theme = useTheme();
  return {
    color: theme.palette.success[500],
  };
});

const WarningIcon = styled(MuiWarningIcon)(() => {
  const theme = useTheme();
  return {
    color: theme.palette.warning[500],
  };
});

const AlertTitle = styled(MuiAlertTitle)(() => {
  const theme = useTheme();
  return {
    color: theme.palette.secondary[900],
    fontWeight: 600,
    fontSize: "16px",
  };
});

const iconMappings = {
  error: <ErrorIcon />,
  info: <InfoIcon />,
  success: <SuccessIcon />,
  warning: <WarningIcon />,
};

export const SEVERITIES = Object.keys(iconMappings);

export interface AlertProps
  extends Pick<
    MuiAlertProps,
    "severity" | "action" | "onClose" | "elevation" | "variant" | "icon"
  > {
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
  <Grid
    container
    direction={direction}
    justifyContent="center"
    alignContent="space-between"
    wrap="nowrap"
  >
    {children}
  </Grid>
);
