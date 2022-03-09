import * as React from "react";
import { Grid } from "@material-ui/core";
import MuiSuccessIcon from "@material-ui/icons/CheckCircle";
import MuiErrorIcon from "@material-ui/icons/Error";
import MuiInfoIcon from "@material-ui/icons/Info";
import MuiWarningIcon from "@material-ui/icons/Warning";
import type { AlertProps as MuiAlertProps } from "@material-ui/lab";
import { Alert as MuiAlert, AlertTitle as MuiAlertTitle } from "@material-ui/lab";

import styled from "../styled";

const backgroundColors = {
  error: "#FDE9E7",
  info: "#EBEDFB",
  success: "#E6F7EB",
  warning: "#FFFDE6",
};

const backgroundGradients = {
  error: `linear-gradient(to right, #DB3615 8px, ${backgroundColors.error} 0%)`,
  info: `linear-gradient(to right, #3548D4 8px, ${backgroundColors.info} 0%)`,
  success: `linear-gradient(to right, #1E942E 8px, ${backgroundColors.success} 0%)`,
  warning: `linear-gradient(to right, #FFCC80 8px, ${backgroundColors.warning} 0%)`,
};

const StyledAlert = styled(MuiAlert)<{ severity: MuiAlertProps["severity"]; $open: boolean }>(
  {
    borderRadius: "8px",
    color: "rgba(13, 16, 48, 0.6)",
    fontSize: "14px",
    overflow: "auto",
    ".MuiAlert-icon": {
      padding: "0",
      margin: "auto",
    },
    ".MuiAlert-message": {
      maxWidth: "calc(100% - 40px)",
      width: "100%",
      padding: "0",
      margin: "auto",
      ".MuiAlertTitle-root": {
        marginBottom: "0",
        color: "#0D1030",
      },
    },
  },
  ({ severity, $open = true }) => ({
    background: $open ? backgroundGradients[severity] : backgroundColors[severity],
    padding: $open ? "16px 16px 20px 24px" : "16px",
    ".MuiAlert-icon": {
      marginRight: $open ? "16px" : null,
    },
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

export interface AlertProps
  extends Pick<MuiAlertProps, "severity" | "action" | "onClose" | "elevation" | "variant"> {
  title?: React.ReactNode;
  collapsible?: boolean;
  defaultOpen?: boolean;
  hover?: boolean;
}

export const Alert: React.FC<AlertProps> = ({
  severity = "info",
  title,
  collapsible = false,
  defaultOpen = false,
  hover = false,
  children,
  ...props
}) => {
  const [open, setOpen] = React.useState<boolean>(defaultOpen);

  return (
    <StyledAlert
      severity={severity}
      iconMapping={iconMappings}
      {...props}
      $open={!collapsible || open}
      onClick={() => !hover && setOpen(!open)}
      onMouseEnter={() => hover && setOpen(true)}
      onMouseLeave={() => hover && setOpen(false)}
    >
      {(!collapsible || open) && title && <AlertTitle>{title}</AlertTitle>}
      {(!collapsible || open) && children}
    </StyledAlert>
  );
};

export interface AlertPanelProps {
  direction?: "row" | "column";
  children: React.ReactElement<AlertProps> | React.ReactElement<AlertProps>[];
}

export const AlertPanel = ({ direction = "column", children }: AlertPanelProps) => (
  <Grid container direction={direction} justify="center" alignContent="space-between" wrap="nowrap">
    {children}
  </Grid>
);
