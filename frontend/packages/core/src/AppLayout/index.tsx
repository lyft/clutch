import React from "react";
import { Outlet } from "react-router-dom";
import { Grid as MuiGrid } from "@mui/material";

import Loadable from "../loading";
import styled from "../styled";
import type { AppConfiguration } from "../Types";

import Drawer from "./drawer";
import { APP_BAR_HEIGHT, Header } from "./header";

const AppGrid = styled(MuiGrid)({
  flex: 1,
});

const ContentGrid = styled(MuiGrid)<{ $isFullScreen: boolean }>(
  {
    flex: 1,
  },
  props => ({
    maxHeight: props.$isFullScreen ? "100vh" : `calc(100vh - ${APP_BAR_HEIGHT})`,
  })
);

const MainContent = styled("div")({ overflowY: "auto", width: "100%" });

interface AppLayoutProps {
  isLoading?: boolean;
  configuration?: AppConfiguration;
  header?: React.ReactElement<any>;
}

const AppLayout: React.FC<AppLayoutProps> = ({
  isLoading = false,
  configuration = {},
  header = null,
}) => {
  return (
    <AppGrid container direction="column" data-testid="app-layout-component">
      {!configuration?.useFullScreenLayout &&
        header &&
        React.cloneElement(header, { ...configuration, ...header.props })}
      {!configuration?.useFullScreenLayout && !header && <Header {...configuration} />}
      <ContentGrid container wrap="nowrap" $isFullScreen={configuration?.useFullScreenLayout}>
        {isLoading ? (
          <Loadable isLoading={isLoading} variant="overlay" />
        ) : (
          <>
            {!configuration?.useFullScreenLayout && <Drawer />}
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
