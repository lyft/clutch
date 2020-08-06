import React from "react";
import { Outlet } from "react-router-dom";

import FeedbackButton from "./feedback";
import Header from "./header";

const AppLayout: React.FC = ({ children }) => {
  return (
    <>
      <Header />
      {children}
      <FeedbackButton />
      <Outlet />
    </>
  );
};

export default AppLayout;
