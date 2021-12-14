import React from "react";
import { matchers } from "@emotion/jest";
import { shallow } from "enzyme";

import { NPSWizard } from "..";

// Add the custom matchers provided by '@emotion/jest'
expect.extend(matchers);

describe("<NPSWizard />", () => {
  describe("basic rendering", () => {
    let wrapper;

    beforeEach(() => {
      wrapper = shallow(<NPSWizard />);
    });

    afterEach(() => {
      wrapper.unmount();
    });

    it("renders", () => {
      expect(wrapper.find(NPSWizard)).toBeDefined();
    });

    it("renders feedback with wizard property", () => {
      expect(wrapper.find("NPSFeedback").props().origin).toBe("WIZARD");
    });

    it("renders the container with a bluish background", () => {
      expect(wrapper).toHaveStyleRule("background", "#F9F9FE");
    });

    it("removes the bluish background after submission", () => {
      wrapper.find("NPSFeedback").prop("onSubmit")(true);
      wrapper.update();
      expect(wrapper).toHaveStyleRule("background", "unset");
    });
  });
});
