import React from "react";
import renderer, { act } from "react-test-renderer";
import { render, screen } from "@testing-library/react";
import type { WorkflowStorageContextProps } from "@clutch-sh/core";
import { client, WorkflowStorageContext } from "@clutch-sh/core";
import { matchers } from "@emotion/jest";
import { mount, shallow } from "enzyme";

import ProjectSelector from "../project-selector";

// Adds the custom matchers provided by '@emotion/jest'
expect.extend(matchers);

const mockedEmptyWorkflowStorageContext: WorkflowStorageContextProps = {
  fromShortLink: false,
  storeData: () => null,
  removeData: () => null,
  retrieveData: () => null,
};

const mockedWorkflowStorageContext: WorkflowStorageContextProps = {
  fromShortLink: false,
  storeData: () => null,
  removeData: () => null,
  retrieveData: () => ({
    "0": {
      clutch: { checked: true, custom: true },
    },
  }),
};

test("Renders", () => {
  const { debug } = render(
    <WorkflowStorageContext.Provider value={mockedEmptyWorkflowStorageContext}>
      <ProjectSelector />
    </WorkflowStorageContext.Provider>
  );
  debug();
  expect(false).toBeTruthy();
});

// describe("<ProjectSelector />", () => {
//   const genWrapper = (storageValue: WorkflowStorageContextProps) =>
//     mount(
//       <WorkflowStorageContext.Provider value={storageValue}>
//         <ProjectSelector />
//       </WorkflowStorageContext.Provider>
//     );

//   describe("basic functionality", () => {
//     it("renders", () => {
//       const wrapper = genWrapper(mockedEmptyWorkflowStorageContext);

//       expect(wrapper.find(<ProjectSelector />)).toBeDefined();
//     });

//     it("renders from storage", () => {
//       const wrapper = genWrapper(mockedWorkflowStorageContext);

//       console.log(wrapper.debug());

//       expect(false).toBeTruthy();
//     });
//   });
// });
