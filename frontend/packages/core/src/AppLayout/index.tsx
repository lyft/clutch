import React from "react";
import { Outlet } from "react-router-dom";
import { Grid as MuiGrid } from "@mui/material";

import type { AppConfiguration } from "../AppProvider";
import Loadable from "../loading";
import styled from "../styled";

import Drawer from "./drawer";
import Header, { APP_BAR_HEIGHT } from "./header";

const AppGrid = styled(MuiGrid)({
  flex: 1,
});

const ContentGrid = styled(MuiGrid)({
  flex: 1,
  maxHeight: `calc(100vh - ${APP_BAR_HEIGHT})`,
});

const MainContent = styled("div")({ overflowY: "auto", width: "100%" });

interface AppLayoutProps {
  isLoading?: boolean;
  configuration?: AppConfiguration;
}

const AppLayout: React.FC<AppLayoutProps> = ({ isLoading = false, configuration = {} }) => {
  return (
    <AppGrid container direction="column">
      <Header {...configuration} />
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
    </AppGrid>
  );
};

export default AppLayout;
