import * as React from "react";
import { DateTimePicker as MuiDateTimePicker } from "@mui/x-date-pickers";
import type { ReactWrapper } from "enzyme";
import { mount } from "enzyme";

import type { DateTimePickerProps } from "../date-time";
import DateTimePicker from "../date-time";

describe("DateTimePicker", () => {
  describe("TextField", () => {
    const onChange = jest.fn();
    let component: ReactWrapper<DateTimePickerProps>;
    beforeAll(() => {
      component = mount(<DateTimePicker value={new Date()} onChange={onChange} />);
    });

    it("has padding", () => {
      const adornedInput = component.find("div.MuiInputBase-adornedEnd");
      expect(adornedInput).toHaveLength(1);
      expect(getComputedStyle(adornedInput.getDOMNode()).getPropertyValue("padding-right")).toBe(
        "14px"
      );
    });

    describe("onChange callback", () => {
      beforeEach(() => {
        onChange.mockReset();
      });

      it("is called when valid value", () => {
        component = mount(<DateTimePicker value={new Date()} onChange={onChange} />);
        const input = component.find("input");
        expect(input).toHaveLength(1);
        input.simulate("change", { target: { value: "11/16/2023 02:55 AM" } });
        expect(onChange).toHaveBeenCalled();
      });

      it("is not called with invalid value", () => {
        component = mount(<DateTimePicker value={new Date()} onChange={onChange} />);
        const input = component.find("input");
        expect(input).toHaveLength(1);
        input.simulate("change", { target: { value: "invalid" } });
        expect(onChange).not.toHaveBeenCalled();
      });
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
      const input = component.find(MuiDateTimePicker);
      expect(input.props().value).toBe(date);
    });

    it("onChange", () => {
      const input = component.find(MuiDateTimePicker);
      expect(input.props().onChange).toBeDefined();
    });

    it("label", () => {
      const outline = component.find(MuiDateTimePicker);
      expect(outline.props().label).toBe("testing");
    });

    it("disabled", () => {
      const outline = component.find(MuiDateTimePicker);
      expect(outline.props().disabled).toBe(false);
    });
  });
});
