import * as React from "react";
import styled from "@emotion/styled";
import type { TabProps as MuiTabProps, TabsProps as MuiTabsProps } from "@material-ui/core";
import { Tab as MuiTab, Tabs as MuiTabs } from "@material-ui/core";

const StyledTab = styled(MuiTab)({
  minWidth: "111px",
  height: "46px",
  padding: "0",
  color: "rgba(13, 16, 48, 0.6)",
  borderBottom: "3px solid #E7E7EA",
  fontSize: "14px",
  opacity: "1",
  "&.Mui-selected": {
    backgroundColor: "unset",
    color: "#3548D4",
    border: "0",
  },
  "&:hover": {
    color: "rgba(13, 16, 48, 0.6)",
    backgroundColor: "#E7E7EA",
    outline: "none",
  },
  "&:focus": {
    color: "#3548D4",
    backgroundColor: "#EBEDFB",
  },
  "&:focus-within": {
    color: "#3548D4",
    backgroundColor: "#EBEDFB",
  },
  "&:active": {
    color: "rgba(13, 16, 48, 0.6)",
    backgroundColor: "rgba(219, 219, 224, 1)",
  },
});

const StyledTabs = styled(MuiTabs)({
  ".MuiTabs-indicator": {
    height: "4px",
    backgroundColor: "#3548D4",
  },
});

export interface TabProps extends Pick<MuiTabProps, "label" | "selected" | "value" | "onClick"> {}

export const Tab = ({ onClick, ...props }: TabProps) => {
  const onClickMiddleware = (e: any) => {
    e.currentTarget.blur();
    if (onClick) {
      onClick(e);
    }
  };
  return <StyledTab color="primary" onClick={onClickMiddleware} {...props} />;
};

export interface TabsProps extends Pick<MuiTabsProps, "value"> {
  children: React.ReactElement<TabProps> | React.ReactElement<TabProps>[];
  // n.b. we explicitly override this prop due to https://github.com/mui-org/material-ui/issues/17454
  onChange:
    | ((event: React.ChangeEvent<{}>, value: any) => void)
    | ((event: React.FormEvent<HTMLButtonElement>) => void);
}

export const Tabs = ({ children, ...props }: TabsProps) => (
  <StyledTabs {...props}>{children}</StyledTabs>
);

export default Tab;
