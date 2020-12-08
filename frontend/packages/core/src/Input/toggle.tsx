import * as React from "react";
import styled from "@emotion/styled";
import type { SwitchProps as MuiSwitchProps } from "@material-ui/core";
import { Switch as MuiSwitch } from "@material-ui/core";

const SwitchContainer = styled(MuiSwitch)({
  ".MuiSwitch-switchBase": {
    ":hover": {
      backgroundColor: "rgba(231, 231, 234, .65)",
    },
    ":focus": {
      backgroundColor: "rgba(226, 226, 230, .65)",
    },
    ":active": {
      backgroundColor: "rgba(219, 219, 224, .65)",
    },
    ".MuiSwitch-thumb": {
      boxShadow: "0px 1px 1px rgba(0, 0, 0, 0.25)",
      color: "#FFFFFF",
    },
  },
  ".MuiSwitch-track": {
    backgroundColor: "#6E7083",
    opacity: 1,
  },
  ".Mui-disabled": {
    ".MuiSwitch-thumb": {
      color: "rgba(248, 248, 249, 1)",
    }
  },
  ".Mui-disabled + .MuiSwitch-track": {
    backgroundColor: "#E7E7EA",
    opacity: 1,
  },
  ".Mui-checked": {
    ":hover": {
      backgroundColor: "rgba(245, 246, 253, .65)",
    },
    ":focus": {
      backgroundColor: "rgba(235, 237, 251, .65)",
    },
    ":active": {
      backgroundColor: "rgba(215, 218, 246, .65)",
    },
    ".MuiSwitch-thumb": {
      color: "#3548D4",
    },
  },
  ".Mui-checked + .MuiSwitch-track": {
    backgroundColor: "#C2C8F2",
    opacity: 1,
  },
  ".Mui-checked.Mui-disabled": {
    ".MuiSwitch-thumb": {
      color: "#E7E7EA",
    },
  },
  ".Mui-checked.Mui-disabled + .MuiSwitch-track": {
    backgroundColor: "#A3A4B0",
    opacity: 1,
  },
});

export interface SwitchProps extends Pick<MuiSwitchProps, "checked" | "disabled" | "onChange"> {};

export const Switch = ({...props}: SwitchProps) => (
  <SwitchContainer color="default" {...props} />
);

export default Switch;
