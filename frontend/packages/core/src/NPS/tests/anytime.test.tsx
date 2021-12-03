import React from "react";
import ChatBubbleOutlineIcon from "@material-ui/icons/ChatBubbleOutline";
import { shallow } from "enzyme";

// import NPSFeedback from "../feedback";
import { NPSAnytime } from "..";

describe("<NPSAnytime />", () => {
  describe("basic rendering", () => {
    it("renders", () => {
      const wrapper = shallow(<NPSAnytime />);
      expect(wrapper.find(NPSAnytime)).toBeDefined();
    });

    it("renders clickable feedback icon", () => {
      const wrapper = shallow(<NPSAnytime />);
      const feedbackIcon = wrapper.find("#anytimeFeedbackIcon");
      expect(feedbackIcon).toBeTruthy();
      expect(feedbackIcon.children().contains(<ChatBubbleOutlineIcon />)).toBeTruthy();
    });
  });

  //   describe("renders feedback window", () => {
  //     const wrapper = shallow(<NPSAnytime />);
  //     const feedbackIcon = wrapper.find("#anytimeFeedbackIcon");

  //     feedbackIcon.props().onClick();

  //     console.log(wrapper.debug());
  //   });
});
