import React from "react";
import { Outlet } from "react-router-dom";
import styled from "@emotion/styled";
import { Box as MuiBox } from "@material-ui/core";

import Drawer from "./drawer";
import FeedbackButton from "./feedback";
import Header from "./header";

const Box = styled(MuiBox)({
  display: "flex",
  ".app-content": {
    width: "calc(100% - 100px)",
  },
  ".divider": {
    minHeight: "64px",
  },
});

const AppLayout: React.FC = () => {
  return (
    <>
      <Box className="app-main">
        <Header />
        <Drawer />
        <FeedbackButton />
        <div className="app-content">
          <div className="divider" />
          <Outlet />
        </div>
      </Box>
    </>
  );
};

export default AppLayout;
