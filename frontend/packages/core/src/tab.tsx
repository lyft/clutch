import * as React from "react";
import styled from "@emotion/styled";
import type { TabProps as MuiTabProps, TabsProps as MuiTabsProps } from "@material-ui/core";
import { Tab as MuiTab } from "@material-ui/core";
import { TabContext, TabList, TabPanel as MuiTabPanel } from "@material-ui/lab";

const StyledTab = styled(MuiTab)({
  minWidth: "111px",
  height: "46px",
  padding: "0",
  color: "rgba(13, 16, 48, 0.6)",
  borderBottom: "3px solid #E7E7EA",
  fontSize: "14px",
  fontWeight: "bold",
  opacity: "1",
  textTransform: "none",
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
    backgroundColor: "#DBDBE0",
  },
});

const StyledTabs = styled(TabList)({
  ".MuiTabs-indicator": {
    height: "4px",
    backgroundColor: "#3548D4",
  },
});

export interface TabProps extends Pick<MuiTabProps, "label" | "selected" | "value" | "onClick"> {
  children?: React.ReactNode;
}

export const Tab = ({ onClick, ...props }: TabProps) => {
  const tabProps = { ...props };
  delete tabProps.children;
  const onClickMiddleware = (e: any) => {
    e.currentTarget.blur();
    if (onClick) {
      onClick(e);
    }
  };
  return <StyledTab color="primary" onClick={onClickMiddleware} {...tabProps} />;
};

const TabPanel = styled(MuiTabPanel)({
  padding: "0",
  maxWidth: "100%",
});

export interface TabsProps extends Pick<MuiTabsProps, "value" | "variant"> {
  children: React.ReactElement<TabProps> | React.ReactElement<TabProps>[];
}

export const Tabs = ({ children, value, variant }: TabsProps) => {
  const [selectedIndex, setSelectedIndex] = React.useState((value || 0).toString());
  const onChangeMiddleware = (_, v: string) => {
    setSelectedIndex(v);
  };

  return (
    <div style={{ width: "100%" }}>
      <TabContext value={selectedIndex}>
        <StyledTabs variant={variant} onChange={onChangeMiddleware}>
          {React.Children.map(children, (child, index) =>
            React.cloneElement(child, { value: index.toString() })
          )}
        </StyledTabs>
        {React.Children.map(children, (tab, index) => (
          <TabPanel value={index.toString()}>{tab.props?.children}</TabPanel>
        ))}
      </TabContext>
    </div>
  );
};
