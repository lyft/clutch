import React from "react";
import { Outlet } from "react-router-dom";
import styled from "@emotion/styled";
import { Grid as MuiGrid } from "@material-ui/core";

import Drawer from "./drawer";
import FeedbackButton from "./feedback";
import Header from "./header";

const AppGrid = styled(MuiGrid)({
  flex: 1,
});

const ContentGrid = styled(MuiGrid)({
  flex: 1,
  overflow: "hidden",
});

const AppLayout: React.FC = () => {
  return (
    <AppGrid container direction="column">
      <Header />
      <ContentGrid container wrap="nowrap">
        <Drawer />
        <Outlet />
      </ContentGrid>
      <FeedbackButton />
    </AppGrid>
  );
};

export default AppLayout;
