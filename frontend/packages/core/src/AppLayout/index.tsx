import React from "react";
import { Outlet } from "react-router-dom";
import styled from "@emotion/styled";
import { Grid as MuiGrid } from "@material-ui/core";

import Loadable from "../loading";

import Drawer from "./drawer";
import FeedbackButton from "./feedback";
import Header from "./header";

const AppGrid = styled(MuiGrid)({
  flex: 1,
});

const ContentGrid = styled(MuiGrid)({
  flex: 1,
  maxHeight: "calc(100vh - 64px)",
});

const MainContent = styled.div({ overflowY: "auto", width: "100%" });

interface AppLayoutProps {
  isLoading?: boolean;
}

const AppLayout: React.FC<AppLayoutProps> = ({ isLoading = false }) => {
  return (
    <AppGrid container direction="column">
      <Header />
      <ContentGrid container wrap="nowrap">
        {isLoading ? (
          <Loadable isLoading={isLoading} variant="overlay" />
        ) : (
          <>
            <Drawer />
            <MainContent>
              <Outlet />
            </MainContent>
          </>
        )}
      </ContentGrid>
      <FeedbackButton />
    </AppGrid>
  );
};

export default AppLayout;
