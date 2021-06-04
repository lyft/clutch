import * as React from "react";
import { shallow } from "enzyme";

import Select from "../select";

describe("Select", () => {
  describe("default option", () => {
    it("has lower bound", () => {
      const component = shallow(
        <Select name="foobar" defaultOption={-1} options={[{ label: "foo" }, { label: "bar" }]} />
      );
      expect(component.find("#foobar")).toHaveLength(1);
      expect(component.find("#foobar-select").props().value).toEqual("foo");
    });

    it("has upper bound", () => {
      const component = shallow(
        <Select name="foobar" defaultOption={2} options={[{ label: "foo" }]} />
      );
      expect(component.find("#foobar")).toHaveLength(1);
      expect(component.find("#foobar-select").props().value).toEqual("foo");
    });
  });
});
