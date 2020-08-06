import React from "react";
import { Outlet } from "react-router-dom";
import { Container } from "@material-ui/core";

import FeedbackButton from "./feedback";
import Header from "./header";

const AppLayout: React.FC = ({ children }) => {
  return (
    <>
      <Header />
      {children}  
      <FeedbackButton />
      <Container maxWidth="xl">
        <Outlet />
      </Container>
    </>
  );
};

export default AppLayout;
