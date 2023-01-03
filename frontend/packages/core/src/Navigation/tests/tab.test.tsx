import React from "react";
import { mount, shallow } from "enzyme";

import { Tab, Tabs } from "../../tab";

describe("Tabs component", () => {
  describe("with a value set", () => {
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

    it("renders children tabs", () => {
      expect(component.find(Tab)).toHaveLength(2);
    });

    it("displays the tab from the specified value", () => {
      expect(component.find("TabContext").prop("value")).toBe("1");
    });

    it("has the mix tab selected", () => {
      // use mount instead of shallow so that we get the other props like
      // `selected`
      const mounted = mount(
        <Tabs value={1}>
          <Tab label="meow" />
          <Tab label="mix" />
        </Tabs>
      );
      expect(mounted.find(Tab).at(1).prop("label")).toBe("mix");
      expect(mounted.find(Tab).at(1).prop("selected")).toBe(true);
    });
  });
});
