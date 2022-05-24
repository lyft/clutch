import React from "react";
import {
  fireEvent,
  render,
  screen,
  waitFor,
  waitForElementToBeRemoved,
  within,
} from "@testing-library/react";
import { act, renderHook } from "@testing-library/react-hooks";
import userEvent from "@testing-library/user-event";

import "@testing-library/jest-dom/extend-expect";
import "@testing-library/jest-dom";

import { DispatchContext, StateContext } from "../helpers";
import ProjectGroup from "../project-group";
import selectorReducer from "../selector-reducer";
import type { State } from "../types";
import { Group } from "../types";

//   screen.debug();
// log entire document to testing-playground
//   screen.logTestingPlaygroundURL();
//   screen.debug();
// log a single element
//   screen.logTestingPlaygroundURL(screen.getByText("test"));

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

// test("Can collapse a Project Group by clicking the header", async () => {
//   const { user } = setup(MockContent());

//   expect(screen.getByLabelText("Collapse Group")).toBeInTheDocument();

//   await user.click(screen.getByLabelText("Test Title Project Group Header"));

//   expect(screen.getByLabelText("Expand Group")).toBeInTheDocument();
// });

// test("Can collapse a Project Group by clicking the icon", async () => {
//   const { user } = setup(MockContent());

//   await user.click(screen.getByLabelText("Collapse Group"));

//   expect(screen.getByLabelText("Expand Group")).toBeInTheDocument();
// });

// test("Will display given projects from the context", () => {
//   mockRender();

//   screen.debug();
//   screen.logTestingPlaygroundURL();

//   expect(false).toBeTruthy();
// });

// const {input} = setup()
// expect(input.value).toBe('') // empty before
// fireEvent.change(input, {target: {value: 'Good Day'}})
// expect(input.value).toBe('') //empty after

test("Can toggle a project", async () => {
  const { user, container } = setup(MockContent());

  const name = "testProject2";
  const inputEl = container.querySelector(`input[name="${name}"]`);

  expect(inputEl).not.toBeChecked();

  await act(async () => user.click(inputEl));

  //   screen.debug();

  //   fireEvent.change(inputEl, { target: { checked: true } });

  //   expect(inputEl).toBeChecked();
});

// test("Can select only a given project", async () => {
//   const { user, container } = setup(MockContent());

//   //   const checkboxes = screen.getAllByRole("checkbox");

//   Object.keys(projects).forEach(projectName => {
//     const inputEl = container.querySelector(`input[name="${projectName}"]`);
//     expect(inputEl).toHaveProperty("checked", projects[projectName].checked);
//   });

//   //   await user.click(screen.getByLabelText("Select Only testProject2 Project"));
//   //   await user.click(
//   //     screen.getByRole("button", { name: "Select Only testProject2 Project", hidden: true })
//   //   );

//   fireEvent(
//     screen.getByRole("button", { name: "Select Only testProject2 Project", hidden: true }),
//     new MouseEvent("click")
//   );

//   //   expect(screen.getByText("Only")).toBeInTheDocument();
//   //   fireEvent.mouseOver(screen.getByText("testProject2"));

//   //   await waitFor(() =>
//   //     screen.getByRole("button", { name: "Select Only testProject2 Project", hidden: true })
//   //   );
//   //   expect(screen.getByText("1st menu item")).toBeInTheDocument();
//   //   expect(screen.getByText("2nd menu item")).toBeInTheDocument();
//   //   expect(screen.getByText("3rd menu item")).toBeInTheDocument();

//   //   screen.debug();

//   Object.keys(projects).forEach(projectName => {
//     const inputEl = container.querySelector(`input[name="${projectName}"]`);
//     expect(inputEl).toHaveProperty("checked", projects[projectName].checked);
//   });

//   //   screen.debug(checkboxes);

//   //   expect(
//   //     within(screen.getByLabelText("Toggle testProject1 Project")).getByRole("checkbox")
//   //   ).toHaveProperty("checked", true);

//   //   expect(
//   //     within(screen.getByLabelText("Toggle testProject2 Project")).getByRole("checkbox")
//   //   ).toHaveProperty("checked", false);

//   //   expect(
//   //     within(screen.getByLabelText("Toggle testProject3 Project")).getByRole("checkbox")
//   //   ).toHaveProperty("checked", true);

//   //   await user.click(screen.getByLabelText("Select Only testProject2 Project"));

//   //   screen.debug();
//   //REMOVE
//   expect(false).toBeTruthy();
// });

// test("Can enable all projects with the switch", async () => {

// });

test("Can remove a project from the list", async () => {
  const { user } = setup(MockContent());

  expect(screen.getByText("testProject2")).toBeInTheDocument();

  await act(async () =>
    user.click(screen.getByRole("button", { name: /remove testproject2 project/i }))
  );

  screen.logTestingPlaygroundURL();

  //   await waitFor(() => {
  //     expect(screen.getByText("testProject2")).not.toBeInTheDocument();
  //   });
  await waitForElementToBeRemoved(() => screen.queryByText("testProject2"));
});

// test("Will display a checked count for the given projects", () => {
//   setup(MockContent());
//   expect(screen.getByText(/2 \/ 3/i)).toBeInTheDocument();
// });

// const checkbox = getByTestId('checkbox-1234').querySelector('input[type="checkbox"]')
// expect(checkbox).toHaveProperty('checked', true)
