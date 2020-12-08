import React from "react";
import { Outlet } from "react-router-dom";

import FeedbackButton from "./feedback";
import Header from "./header";

const AppLayout: React.FC = ({ children }) => {
  return (
    <>
      <div style={{display: "flex"}}>
      <Header />
        <div style={{minHeight: "64px"}}/>
        {children}
      </div>
      <FeedbackButton />
      <Outlet />
    </>
  );
};

export default AppLayout;
