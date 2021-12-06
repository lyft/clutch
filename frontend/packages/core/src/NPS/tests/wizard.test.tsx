import React from "react";
import { shallow } from "enzyme";

import NPSFeedback from "../feedback";
import { NPSWizard } from "..";

describe("<NPSWizard />", () => {
  describe("basic rendering", () => {
    let wrapper;

    beforeEach(() => {
      wrapper = shallow(<NPSWizard />);
    });

    it("renders", () => {
      expect(wrapper.find(NPSWizard)).toBeDefined();
    });

    it("renders feedback with wizard property", () => {
      expect(wrapper.contains(<NPSFeedback origin="WIZARD" />)).toEqual(true);
    });
  });
});
