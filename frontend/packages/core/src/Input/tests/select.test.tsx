import * as React from "react";
import { shallow } from "enzyme";

import { Select } from "../select";

describe("Select", () => {
  describe("default option", () => {
    it("has lower bound", () => {
      const component = shallow(
        <Select name="foobar" defaultOption={-1} options={[{ label: "foo" }, { label: "bar" }]} />
      );
      expect(component.find("#foobar")).toHaveLength(1);
      expect(component.find("#foobar-select").props().value).toStrictEqual(["foo"]);
    });

    it("has upper bound", () => {
      const component = shallow(
        <Select name="foobar" defaultOption={2} options={[{ label: "foo" }]} />
      );
      expect(component.find("#foobar")).toHaveLength(1);
      expect(component.find("#foobar-select").props().value).toStrictEqual(["foo"]);
    });
  });

  describe("multiple values", () => {
    it("allows multiple", () => {
      const component = shallow(
        <Select
          name="foobar"
          multiple
          defaultOption={2}
          options={[{ label: "foo" }, { label: "bar" }]}
        />
      );
      expect(component.find("#foobar")).toHaveLength(1);
      component.find("#foobar-select").simulate("change", { target: { value: ["foo", "bar"] } });
      expect(component.find("#foobar-select").props().value).toStrictEqual(["foo", "bar"]);
    });
  });
});
