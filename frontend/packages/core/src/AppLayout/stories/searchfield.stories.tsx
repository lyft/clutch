import * as React from "react";
import { MemoryRouter } from "react-router";
import type { Meta } from "@storybook/react";

import { ApplicationContext } from "../../Contexts/app-context";
import SearchField from "../search";

export default {
  title: "Core/AppLayout/Search Field",
  component: SearchField,
  decorators: [
    () => (
      <MemoryRouter>
        <SearchField />
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
  parameters: {
    backgrounds: {
      default: "clutch",
      values: [{ name: "clutch", value: "#131C5F" }],
    },
  },
} as Meta;

export const Primary: React.FC<{}> = () => <SearchField />;
