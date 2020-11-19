import React from "react";
import { Link } from "react-router-dom";
import styled from "@emotion/styled";
import { AppBar as MuiAppBar, Box, Grid, IconButton, Toolbar, Typography } from "@material-ui/core";
import MenuIcon from "@material-ui/icons/Menu";

import Drawer from "./drawer";
import Logo from "./logo";
import Notifications from "./notifcations";
import SearchField from "./search";
import { UserInformation } from "./user";

const AppBar = styled(MuiAppBar)`
  min-width: fit-content;
  background: linear-gradient(90deg, #38106b 4.58%, #131c5f 89.31%);
  min-width: fit-content;
`;

const MenuButton = styled(IconButton)`
  padding: 12px;
  margin-left: -12px;
`;

const Title = styled(Typography)`
  margin-right: 25px;
  font-weight: bold;
  font-size: 20px;
  line-height: 24px;
  color: #ffffff;
  opacity: 0.87;
`;

const Header: React.FC = () => {
  const [drawerOpen, setDrawerOpen] = React.useState(false);

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
      <AppBar position="static" color="secondary">
        <Toolbar>
          <MenuButton onClick={openDrawer} edge="start" color="primary" data-qa="menuBtn">
            <MenuIcon />
          </MenuButton>
          <Link to="/">
            <Logo />
          </Link>
          <Title>clutch</Title>
          <Box />
          <SearchField />
          <Box />
          <Grid container alignItems="center" justify="flex-end">
            <Notifications />
            <UserInformation />
          </Grid>
        </Toolbar>
      </AppBar>
      <Drawer open={drawerOpen} onClose={onDrawerClose} />
    </>
  );
};

export default Header;
