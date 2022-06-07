import React from "react";
import { fireEvent, render, screen } from "@testing-library/react";
import { renderHook } from "@testing-library/react-hooks";
import userEvent from "@testing-library/user-event";

import "@testing-library/jest-dom/extend-expect";
import "@testing-library/jest-dom";

import { DispatchContext, StateContext } from "../helpers";
import ProjectGroup from "../project-group";
import selectorReducer from "../selector-reducer";
import type { State } from "../types";
import { Group } from "../types";

const projects = {
  testProject1: { checked: true },
  testProject2: { checked: false, custom: true },
  testProject3: { checked: true, custom: true },
};

const stateOptions = {
  projectData: {},
  projectErrors: [],
  loading: false,
  error: undefined,
};

const emptyState: State = {
  [Group.PROJECTS]: {},
  [Group.UPSTREAM]: {},
  [Group.DOWNSTREAM]: {},
  ...stateOptions,
};

const mockState: State = {
  [Group.PROJECTS]: projects,
  [Group.UPSTREAM]: {},
  [Group.DOWNSTREAM]: {},
  ...stateOptions,
};

const setup = jsx => ({
  user: userEvent.setup(),
  ...render(jsx),
});

const MockContent = (initState = mockState) => {
  // eslint-disable-next-line react-hooks/rules-of-hooks
  const { result } = renderHook(() => React.useReducer(selectorReducer, initState));

  const [state, dispatch] = result.current;

  return (
    <DispatchContext.Provider value={dispatch}>
      <StateContext.Provider value={state}>
        <ProjectGroup title="Test Title" group={Group.PROJECTS} />
      </StateContext.Provider>
    </DispatchContext.Provider>
  );
};

test("Renders with a given Title", () => {
  setup(MockContent());
  expect(screen.getByText("Test Title")).toBeInTheDocument();
});

test("Renders a message when no projects are defined", () => {
  setup(MockContent(emptyState));
  expect(screen.getByText("No projects in this group yet.")).toBeInTheDocument();
});

test("Can collapse a Project Group by clicking the header", async () => {
  const { user } = setup(MockContent());

  expect(screen.getByLabelText("Collapse Group")).toBeInTheDocument();

  await user.click(screen.getByLabelText("Test Title Project Group Header"));

  expect(screen.getByLabelText("Expand Group")).toBeInTheDocument();
});

test("Can collapse a Project Group by clicking the icon", async () => {
  const { user } = setup(MockContent());

  await user.click(screen.getByLabelText("Collapse Group"));

  expect(screen.getByLabelText("Expand Group")).toBeInTheDocument();
});

test("Will display a checked count for the given projects", () => {
  setup(MockContent());

  expect(screen.getByText(/2 \/ 3/i)).toBeInTheDocument();
});

test("Can toggle a project", async () => {
  setup(MockContent());

  const checkboxContainer = screen.getByLabelText("Toggle testProject2 Project");

  const checkbox = checkboxContainer.querySelector("input[type='checkbox']");

  expect(checkbox).not.toBeChecked();

  fireEvent.change(checkbox, { target: { checked: true } });

  expect(checkbox).toBeChecked();
});

// TODO (jslaughter): Finish test cases
/* eslint-disable jest/no-disabled-tests */
test.skip("Can select only a given project", async () => {});
test.skip("Can enable all projects with the switch", async () => {});
test.skip("Can remove a project from the list", async () => {});
/* eslint-enable jest/no-disabled-tests */
