import React from "react";
import { shallow } from "enzyme";

import { Tab, Tabs } from "../tab";

describe("Tabs component aka TabGroup", () => {
  describe("basic rendering", () => {
    let component;
    beforeEach(() => {
      component = shallow(
        <Tabs value={1}>
          <Tab label="meow" />
          <Tab label="mix" />
        </Tabs>
      );
    });

    it("renders", () => {
      expect(component.find(Tabs)).toBeDefined();
    });

    it("renders 2 Tabs", () => {
      expect(component.find(Tab)).toHaveLength(2);
    });

    it("has the second tab selected (index 1)", () => {
      expect(component.find("TabContext").prop("value")).toBe("1");
    });
  });
});
