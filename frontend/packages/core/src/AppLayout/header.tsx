import React from "react";
import { Link } from "react-router-dom";
import {
  AppBar as MuiAppBar,
  Box,
  Divider as MuiDivider,
  IconButton,
  Toolbar,
  Typography,
} from "@material-ui/core";
import MenuIcon from "@material-ui/icons/Menu";
import styled from "styled-components";

import Drawer from "./drawer";
import Logo from "./logo";
import SearchField from "./search";
import { UserInformation } from "./user";

const AppBar = styled(MuiAppBar)`
  min-width: fit-content;
`;

const Title = styled(Typography)`
  margin-right: 25px;
  font-weight: bolder;
`;

const Divider = styled(MuiDivider)`
  ${({ theme }) => `
  background-color: ${theme.palette.primary.main};
  margin: 16px 8px;
  `}
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
      <AppBar position="static" color="secondary" style={{ minWidth: "fit-content" }}>
        <Toolbar>
          <IconButton onClick={openDrawer} edge="start" color="primary" data-qa="menuBtn">
            <MenuIcon />
          </IconButton>
          <Link to="/">
            <Logo />
          </Link>
          <Divider orientation="vertical" flexItem />
          <Title variant="h5">clutch</Title>
          <Box />
          <SearchField />
          <Box />
          <UserInformation />
        </Toolbar>
      </AppBar>
      <Drawer open={drawerOpen} onClose={onDrawerClose} />
    </>
  );
};

export default Header;
