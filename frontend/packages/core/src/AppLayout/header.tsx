import React from "react";
import { Link } from "react-router-dom";
import styled from "@emotion/styled";
import { AppBar as MuiAppBar, Box, Grid, Theme, Toolbar, Typography } from "@mui/material";

import type { AppConfiguration } from "../AppProvider";
import { FeatureOn, SimpleFeatureFlag } from "../flags";
import { NPSHeader } from "../NPS";

import Logo from "./logo";
import Notifications from "./notifications";
import SearchField from "./search";
import ShortLinker from "./shortLinker";
import ThemeSwitcher from "./theme-switcher";
import { UserInformation } from "./user";

export const APP_BAR_HEIGHT = "64px";

/**
 * Properties used to allow Storybook examples to override featureflag settings
 */
interface HeaderProps extends AppConfiguration {
  /**
   * Will enable the NPS feedback component in the header
   */
  enableNPS?: boolean;
  /**
   * Will enable the workflow search component in the header
   */
  search?: boolean;
  /**
   * Will enable the NPS feedback component in the header
   */
  feedback?: boolean;
  /**
   * Will enable the shortlinks component in the header
   */
  shortLinks?: boolean;
  /**
   * Will enable the notifications component in the header
   */
  notifications?: boolean;
  /**
   * Will enable the user information component in the header
   */
  userInfo?: boolean;
}

const AppBar = styled(MuiAppBar)(({ theme }: { theme: Theme }) => ({
  minWidth: "fit-content",
  background: theme.palette.headerGradient,
  zIndex: 1201,
  height: APP_BAR_HEIGHT,
}));

// Since the AppBar is fixed we need a div to take up its height in order to push other content down.
const ClearAppBar = styled.div({ height: APP_BAR_HEIGHT });

const Title = styled(Typography)(({ theme }: { theme: Theme }) => ({
  margin: "12px 0px 12px 8px",
  fontWeight: "bold",
  fontSize: "30px",
  paddingLeft: "5px",
  color: theme.palette.common.white,
}));

const StyledLogo = styled("img")({
  width: "48px",
  height: "48px",
  padding: "1px",
  verticalAlign: "middle",
});

const Header: React.FC<HeaderProps> = ({
  title = "clutch",
  logo = <Logo />,
  enableNPS = false,
  search = true,
  feedback = true,
  shortLinks = true,
  notifications = false,
  userInfo = true,
  children = null,
}) => {
  return (
    <>
      <AppBar position="fixed" elevation={0}>
        <Toolbar>
          <Link to="/">{typeof logo === "string" ? <StyledLogo src={logo} /> : logo}</Link>
          <Title>{title}</Title>
          <Grid container alignItems="center" justifyContent="flex-end">
            {search && (
              <Box>
                <SearchField />
              </Box>
            )}
            {shortLinks && (
              <SimpleFeatureFlag feature="shortLinks">
                <FeatureOn>
                  <ShortLinker />
                </FeatureOn>
              </SimpleFeatureFlag>
            )}
            {notifications && <Notifications />}
            {feedback && (
              <>
                {enableNPS ? (
                  <NPSHeader />
                ) : (
                  <SimpleFeatureFlag feature="npsHeader">
                    <FeatureOn>
                      <NPSHeader />
                    </FeatureOn>
                  </SimpleFeatureFlag>
                )}
              </>
            )}
            {children && children}
            {userInfo && (
              <UserInformation>
                <ThemeSwitcher />
              </UserInformation>
            )}
          </Grid>
        </Toolbar>
      </AppBar>
      <ClearAppBar />
    </>
  );
};

export default Header;
