import React from "react";
import { matchers } from "@emotion/jest";
import ChatBubbleOutlineIcon from "@material-ui/icons/ChatBubbleOutline";
import { shallow } from "enzyme";

import * as ApplicationContext from "../../Contexts/app-context";
import contextValues from "../../Contexts/tests/testContext";
import { NPSHeader } from "..";

// Add the custom matchers provided by '@emotion/jest'
expect.extend(matchers);

describe("<NPSHeader />", () => {
  describe("basic rendering", () => {
    let wrapper;

    beforeEach(() => {
      jest.spyOn(ApplicationContext, "useAppContext").mockReturnValue(contextValues);
      jest.useFakeTimers();
      wrapper = shallow(<NPSHeader />);
    });

    afterEach(() => {
      wrapper.unmount();
    });

    it("renders", () => {
      expect(wrapper.find(NPSHeader)).toBeDefined();
    });

    it("renders clickable feedback icon", () => {
      const feedbackIcon = wrapper.find("#anytimeFeedbackIcon");
      expect(feedbackIcon).toBeTruthy();
      expect(feedbackIcon.children().contains(<ChatBubbleOutlineIcon />)).toBeTruthy();
    });

    it("opens a popper on click of feedback icon", () => {
      const feedbackIcon = wrapper.find("#anytimeFeedbackIcon");

      expect(wrapper.find("Styled(Component)").at(1).prop("open")).toBeFalsy();

      feedbackIcon.prop("onClick")();

      wrapper.update();

      expect(wrapper.find("Styled(Component)").at(1).prop("open")).toBeTruthy();
    });
  });
});
