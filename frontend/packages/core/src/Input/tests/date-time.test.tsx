import * as React from "react";
import { mount } from "enzyme";

import DateTimePicker from "../date-time";

describe("DateTimePicker", () => {
  describe("TextField", () => {
    let component;
    beforeAll(() => {
      component = mount(<DateTimePicker value={new Date()} onChange={() => {}} />);
    });

    it("has padding", () => {
      const adornedInput = component.find("div.MuiInputBase-adornedEnd");
      expect(adornedInput).toHaveLength(1);
      expect(getComputedStyle(adornedInput.getDOMNode()).getPropertyValue("padding-right")).toBe(
        "14px"
      );
    });
  });

  describe("proxies prop", () => {
    let date: Date;
    let onChange: () => void;
    let component;
    beforeAll(() => {
      date = new Date();
      onChange = () => {};
      component = mount(
        <DateTimePicker value={date} onChange={onChange} label="testing" disabled={false} />
      );
    });

    it("value", () => {
      expect(component).toHaveLength(1);
      const input = component.find("ForwardRef(DateTimePicker)");
      expect(input.props().value).toBe(date);
    });

    it("onChange", () => {
      const input = component.find("ForwardRef(DateTimePicker)");
      expect(input.props().onChange).toBe(onChange);
    });

    it("label", () => {
      const outline = component.find("ForwardRef(DateTimePicker)");
      expect(outline.props().label).toBe("testing");
    });

    it("disabled", () => {
      const outline = component.find("ForwardRef(DateTimePicker)");
      expect(outline.props().disabled).toBe(false);
    });
  });
});
