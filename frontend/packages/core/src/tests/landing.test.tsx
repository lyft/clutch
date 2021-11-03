import React from "react";
import { MemoryRouter } from "react-router-dom";
import { shallow, ShallowWrapper } from "enzyme";

import type { Workflow } from "../AppProvider/workflow";
import * as appContext from "../Contexts/app-context";
import Landing from "../landing";

describe("Landing component", () => {
  let component: ShallowWrapper;

  // const trendingWorkflows: Workflow[] = [
  //   {
  //     developer: {
  //       name: "Test",
  //       contactUrl: "mailto:example@example.com",
  //     },
  //     displayName: "EX",
  //     group: "EC",
  //     path: "ex",
  //     routes: [
  //       {
  //         path: "test1",
  //         displayName: "Example Route 1",
  //         description: "Example Description 1",
  //         component: null,
  //         trending: false,
  //       },
  //       {
  //         path: "test2",
  //         displayName: "Example Route 2",
  //         description: "Example Description 2",
  //         component: null,
  //         trending: true,
  //       },
  //     ],
  //   },
  // ];

  beforeAll(() => {
    jest.spyOn(appContext, "useAppContext").mockReturnValue({ workflows: [] });
    component = shallow(
      <MemoryRouter>
        <Landing />
      </MemoryRouter>
    );
  });

  it("renders correctly", () => {
    console.log(component.debug());
    expect(component).toMatchSnapshot();
  });

  // it("renders trending workflows", () => {
  //   jest.spyOn(appContext, "useAppContext").mockReturnValue({ workflows });


  // });

  // it("will not prepend the group if it mathes the route", () => {

  // })
});
