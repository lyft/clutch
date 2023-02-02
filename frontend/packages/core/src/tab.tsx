import * as React from "react";
import styled from "@emotion/styled";
import { TabContext, TabList, TabPanel as MuiTabPanel } from "@mui/lab";
import type { TabProps as MuiTabProps, TabsProps as MuiTabsProps } from "@mui/material";
import { Tab as MuiTab } from "@mui/material";

const StyledTab = styled(MuiTab)({
  minWidth: "111px",
  height: "46px",
  padding: "12px 32px",
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
  ".MuiTab-wrapper": {
    margin: "auto",
  },
});

const StyledTabs = styled(TabList)({
  ".MuiTabs-indicator": {
    height: "4px",
    backgroundColor: "#3548D4",
  },
});

export interface TabProps extends Pick<MuiTabProps, "label" | "value" | "onClick"> {
  children?: React.ReactNode;
  startAdornment?: React.ReactNode;
  selected?: boolean;
}

export const Tab = ({ children, onClick, label, startAdornment, ...props }: TabProps) => {
  const tabProps = { ...props };

  const onClickMiddleware = (e: any) => {
    e.currentTarget.blur();
    if (onClick) {
      onClick(e);
    }
  };
  let finalLabel = label;
  if (startAdornment !== undefined) {
    finalLabel = (
      <div style={{ display: "flex", alignItems: "center" }}>
        <span style={{ display: "inherit", marginRight: "7px" }}>{startAdornment}</span>
        {label}
      </div>
    );
  }
  return <StyledTab color="primary" onClick={onClickMiddleware} label={finalLabel} {...tabProps} />;
};

const TabPanel = styled(MuiTabPanel)({
  padding: "0",
  maxWidth: "100%",
});

export interface TabsProps extends Pick<MuiTabsProps, "value" | "variant" | "centered"> {
  children: React.ReactElement<TabProps> | React.ReactElement<TabProps>[];
  // To allow for callback functionality, use `onChange`
  // Note the value is referring to the index, like "0", "1", "2", etc.
  onChange?: (value: string) => void;
}

export const Tabs = ({ children, value, variant, onChange, ...props }: TabsProps) => {
  const [selectedIndex, setSelectedIndex] = React.useState((value || 0).toString());
  const onChangeMiddleware = (_, v: string) => {
    setSelectedIndex(v);
    if (onChange) {
      onChange(v);
    }
  };

  return (
    <div style={{ width: "100%" }}>
      <TabContext value={selectedIndex}>
        <StyledTabs
          data-testid="styled-tabs"
          variant={variant}
          onChange={onChangeMiddleware}
          {...props}
        >
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
