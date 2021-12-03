import React from "react";
import { shallow } from "enzyme";

import NPSFeedback from "../feedback";

describe("<NPSFeedback />", () => {
  describe("basic rendering", () => {
    it("renders feedback with given origin", () => {
      const wizardWrapper = shallow(<NPSFeedback origin="WIZARD" />);
      expect(wizardWrapper).toBeTruthy();

      const anytimeWrapper = shallow(<NPSFeedback origin="ANYTIME" />);
      expect(anytimeWrapper).toBeTruthy();

      const unspecifiedWrapper = shallow(<NPSFeedback origin="ORIGIN_UNSPECIFIED" />);
      expect(unspecifiedWrapper).toBeTruthy();
    });
  });
});
