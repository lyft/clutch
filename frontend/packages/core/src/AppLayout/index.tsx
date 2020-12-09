import { Box } from "@material-ui/core";
import React from "react";
import { Outlet } from "react-router-dom";

import FeedbackButton from "./feedback";
import Header from "./header";
import Drawer from "./drawer";

const AppLayout: React.FC = ({ children }) => {
  return (
    <>
      <Box style={{display: "flex"}}>
      <Header />
      <Drawer/>
        <div style={{minHeight: "64px"}}/>
        {children}
      </Box>
      <FeedbackButton />
      <Outlet />
    </>
  );
};

export default AppLayout;
