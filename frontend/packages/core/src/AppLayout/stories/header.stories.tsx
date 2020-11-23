import * as React from "react";
import { MemoryRouter } from "react-router";
import type { Meta } from "@storybook/react";

import { ApplicationContext } from "../../Contexts/app-context";
import Header from "../header";

export default {
  title: "Core/AppLayout/Header",
  component: Header,
  decorators: [
    Header => (
      <MemoryRouter>
        <Header />
      </MemoryRouter>
    ),
    StoryFn => {
      return (
        <ApplicationContext.Provider value={{ workflows: [] }}>
          <StoryFn />
        </ApplicationContext.Provider>
      );
    },
  ],
} as Meta;

export const Primary: React.FC<{}> = () => <Header />;
