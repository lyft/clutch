import React from "react";
import { BrowserRouter as Router } from "react-router-dom";
import type { Meta } from "@storybook/react";

import { ApplicationContext } from "../Contexts/app-context";
import Landing from "../landing";

export default {
  title: "Core/Landing",
  decorators: [
    story => (
      // eslint-disable-next-line react/jsx-no-constructed-context-values
      <ApplicationContext.Provider value={{ workflows: [] }}>
        <Router>{story()}</Router>
      </ApplicationContext.Provider>
    ),
  ],
  component: Landing,
} as Meta;

const Template = () => <Landing />;

export const Primary = Template.bind({});
