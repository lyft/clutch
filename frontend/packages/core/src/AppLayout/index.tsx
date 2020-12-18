import React from "react";
import { Outlet } from "react-router-dom";
import { Box, Grid } from "@material-ui/core";

import Drawer from "./drawer";
import FeedbackButton from "./feedback";
import Header from "./header";

const AppLayout: React.FC = () => {
  return (
    <Grid container direction="column" style={{flex: 1}}>
      <Header />
      <Grid container style={{flex: 1}}>
        <Drawer />
        <Outlet />
      </Grid>
      <FeedbackButton />
    </Grid>
  );
};

export default AppLayout;
