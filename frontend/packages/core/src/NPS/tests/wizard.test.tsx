import React from "react";
import { shallow } from "enzyme";

import NPSFeedback from "../feedback";
import { NPSWizard } from "..";

describe("<NPSWizard />", () => {
  describe("basic rendering", () => {
    it("renders", () => {
      const wrapper = shallow(<NPSWizard />);
      expect(wrapper.find(NPSWizard)).toBeDefined();
    });

    it("renders feedback with wizard property", () => {
      const wrapper = shallow(<NPSWizard />);
      expect(wrapper.contains(<NPSFeedback origin="WIZARD" />)).toEqual(true);
    });
  });
});
