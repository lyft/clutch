import React from "react";
import { Link } from "react-router-dom";
import styled from "@emotion/styled";
import { AppBar as MuiAppBar, Box, Grid, IconButton, Toolbar, Typography } from "@material-ui/core";
import MenuIcon from "@material-ui/icons/Menu";

import Drawer from "./drawer";
import Logo from "./logo";
import Notifications from "./notifications";
import SearchField from "./search";
import { UserInformation } from "./user";

// TODO (sperry): make header responsive for small devices
const AppBar = styled(MuiAppBar)({
  minWidth: "fit-content",
  background: "linear-gradient(90deg, #38106b 4.58%, #131c5f 89.31%)",
});

const MenuButton = styled(IconButton)({
  padding: "12px",
  marginLeft: "-12px",
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

  const handleKeyPress = (event: KeyboardEvent) => {
    // @ts-ignore
    if (event.key === "." && event.target?.nodeName !== "INPUT") {
      setDrawerOpen(true);
    } else if (event.key === "Escape") {
      setDrawerOpen(false);
    }
  };

  React.useLayoutEffect(() => {
    window.addEventListener("keydown", handleKeyPress);
  }, []);

  const openDrawer = () => {
    setDrawerOpen(true);
  };

  const onDrawerClose = () => {
    setDrawerOpen(false);
  };

  return (
    <>
      <AppBar position="static" elevation={0}>
        <Toolbar>
          <MenuButton onClick={openDrawer} edge="start" color="primary" data-qa="menuBtn">
            <MenuIcon />
          </MenuButton>
          <Link to="/">
            <Logo />
          </Link>
          <Title>clutch</Title>
          <Grid container alignItems="center" justify="flex-end">
            <Box>
              <SearchField />
            </Box>
            {showNotifications ? <Notifications /> : null}
            <UserInformation />
          </Grid>
        </Toolbar>
      </AppBar>
      <Drawer open={drawerOpen} onClose={onDrawerClose} />
    </>
  );
};

export default Header;
