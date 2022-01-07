import React from "react";
import { matchers } from "@emotion/jest";
import ChatBubbleOutlineIcon from "@material-ui/icons/ChatBubbleOutline";
import { mount, shallow } from "enzyme";

import { NPSHeader } from "..";

// Add the custom matchers provided by '@emotion/jest'
expect.extend(matchers);

describe("<NPSHeader />", () => {
  describe("basic rendering", () => {
    let wrapper;

    beforeEach(() => {
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

      feedbackIcon.prop("onClick")(null);

      wrapper.update();

      expect(wrapper.find("Styled(Component)").at(1).prop("open")).toBeTruthy();
    });

    it("closes the popper on click outside", () => {});

    // it("renders anytime feedback inside of popper", () => {
    //   const mounted = mount(<NPSAnytime />);
    //   console.log(mounted.debug());
    //   const feedbackIcon = mounted.find("#anytimeFeedbackIcon");
    //   feedbackIcon.prop("onClick")(null);

    //   mounted.update();

    //   console.log(mounted.debug());
    // });
  });
});
