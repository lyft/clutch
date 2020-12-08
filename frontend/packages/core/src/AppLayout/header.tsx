import React from "react";
import { Link } from "react-router-dom";
import styled from "@emotion/styled";
import { AppBar as MuiAppBar, Box, Grid, Toolbar, Typography } from "@material-ui/core";

import Drawer from "./drawer";
import Logo from "./logo";
import Notifications from "./notifications";
import SearchField from "./search";
import { UserInformation } from "./user";

// TODO (sperry): make header responsive for small devices
const AppBar = styled(MuiAppBar)({
  minWidth: "fit-content",
  background: "linear-gradient(90deg, #38106b 4.58%, #131c5f 89.31%)",
  zIndex: 1201,
});

const Title = styled(Typography)({
  margin: "12px 0px 12px 8px",
  fontWeight: "bold",
  fontSize: "20px",
  color: "rgba(255, 255, 255, 0.87)",
});

const Header: React.FC = () => {
  const [drawerOpen, setDrawerOpen] = React.useState(false);
  const showNotifications = false;
  const onDrawerClose = () => {
    setDrawerOpen(false);
  };
  return (
    <>
      <AppBar position="static" elevation={0}>
        <Toolbar>
          <Link to="/">
            <Logo />
          </Link>
          <Title>clutch</Title>
          <Grid container alignItems="center" justify="flex-end">
            <Box>
              <SearchField />
            </Box>
            {showNotifications && <Notifications />}
            <UserInformation />
          </Grid>
        </Toolbar>
      </AppBar>
      <Drawer onClose={onDrawerClose} />
    </>
  );
};

export default Header;
