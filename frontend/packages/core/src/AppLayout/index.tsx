import React from "react";
import { Outlet } from "react-router-dom";
import styled from "@emotion/styled";
import { Grid as MuiGrid } from "@material-ui/core";

import Drawer from "./drawer";
import FeedbackButton from "./feedback";
import Header from "./header";

const Grid = styled(MuiGrid)({
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
      <Grid direction="column" className="app-main">
        <Header />
        <Drawer />
        <FeedbackButton />
        <div className="app-content">
          <div className="divider" />
          <Outlet />
        </div>
      </Grid>
    </>
  );
};

export default AppLayout;
