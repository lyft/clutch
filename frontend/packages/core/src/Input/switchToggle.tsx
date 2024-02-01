/*
  TODO from andrewo: Rename this file to switch.tsx or the component itself to something like switchToggle.
  There's an issue with storybook with components whose names are also keywords:
  https://github.com/storybookjs/storybook/issues/11980
*/

import * as React from "react";
import styled from "@emotion/styled";
import type { SwitchProps as MuiSwitchProps, Theme } from "@mui/material";
import { alpha, Switch as MuiSwitch } from "@mui/material";

const SwitchContainer = styled(MuiSwitch)(({ theme }: { theme: Theme }) => ({
  ".MuiSwitch-switchBase": {
    ":hover": {
      backgroundColor: alpha(theme.palette.secondary[900], 0.1),
    },
    ":focus": {
      backgroundColor: alpha(theme.palette.secondary[900], 0.12),
    },
    ":active": {
      backgroundColor: alpha(theme.palette.secondary[900], 0.15),
    },
    ".MuiSwitch-thumb": {
      boxShadow: `0px 1px 1px ${alpha(
        theme.palette.getContrastText(theme.palette.contrastColor),
        0.25
      )}`,
      color: theme.palette.contrastColor,
    },
  },
  ".MuiSwitch-track": {
    backgroundColor: theme.palette.secondary[500],
    opacity: 1,
  },
  ".Mui-disabled": {
    ".MuiSwitch-thumb": {
      color: theme.palette.secondary[50],
    },
  },
  ".Mui-disabled + .MuiSwitch-track": {
    backgroundColor: theme.palette.secondary[200],
    opacity: 1,
  },
  ".Mui-checked": {
    ":hover": {
      backgroundColor: alpha(theme.palette.primary[600], 0.05),
    },
    ":focus": {
      backgroundColor: alpha(theme.palette.primary[600], 0.1),
    },
    ":active": {
      backgroundColor: alpha(theme.palette.primary[600], 0.2),
    },
    ".MuiSwitch-thumb": {
      color: theme.palette.primary[600],
    },
  },
  ".Mui-checked + .MuiSwitch-track": {
    backgroundColor: theme.palette.primary[300],
    opacity: 1,
  },
  ".Mui-checked.Mui-disabled": {
    ".MuiSwitch-thumb": {
      color: theme.palette.secondary[200],
    },
  },
  ".Mui-checked.Mui-disabled + .MuiSwitch-track": {
    backgroundColor: theme.palette.secondary[400],
    opacity: 1,
  },
}));

export interface SwitchProps extends Pick<MuiSwitchProps, "checked" | "disabled" | "onChange"> {}

const Switch = ({ ...props }: SwitchProps) => <SwitchContainer color="default" {...props} />;

export default Switch;
