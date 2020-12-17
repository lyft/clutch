import React from "react";
import { Outlet } from "react-router-dom";
import { Box, Grid } from "@material-ui/core";

import Drawer from "./drawer";
import FeedbackButton from "./feedback";
import Header from "./header";

const AppLayout: React.FC = () => {
  return (
    <>
      <Grid direction="column">
        <Header />
        <Box display="flex">
          <Drawer />
          <Outlet />
        </Box>
        <FeedbackButton />
      </Grid>
    </>
  );
};

export default AppLayout;
