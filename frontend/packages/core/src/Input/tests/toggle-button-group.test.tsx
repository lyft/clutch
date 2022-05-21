import * as React from "react";
import { ToggleButton } from "@material-ui/lab";
import { shallow } from "enzyme";

import ToggleButtonGroup from "../toggle-button-group";

// Note interaction tests will happen when there is a parent component that
// can maintain state and select a default
describe("ToggleButtonGroup", () => {
  let multipleFalseWrapper;
  let multipleTrueWrapper;
  beforeEach(() => {
    multipleFalseWrapper = shallow(
      <ToggleButtonGroup value="foo" onChange={() => {}} size="large">
        <ToggleButton value="foo">foo</ToggleButton>
        <ToggleButton value="Ingress">Ingress</ToggleButton>
        <ToggleButton value="Egress">Egress</ToggleButton>
      </ToggleButtonGroup>
    );

    multipleTrueWrapper = shallow(
      <ToggleButtonGroup
        multiple
        value={["foo", "Ingress"]}
        onChange={() => {}}
        orientation="vertical"
      >
        <ToggleButton value="foo">foo</ToggleButton>
        <ToggleButton value="Ingress">Ingress</ToggleButton>
      </ToggleButtonGroup>
    );
  });
  const searchStringForComponent = "Styled(WithStyles(ForwardRef(ToggleButton)))";

  describe("rendering", () => {
    it("renders", () => {
      expect(multipleFalseWrapper.find(ToggleButtonGroup)).toBeDefined();
    });

    it("renders with 3 toggles", () => {
      expect(multipleFalseWrapper.find(ToggleButton)).toHaveLength(3);
    });

    it("renders with 2 toggles", () => {
      expect(multipleTrueWrapper.find(ToggleButton)).toHaveLength(2);
    });

    it('renders with the starting value being "foo"', () => {
      expect(multipleFalseWrapper.find(searchStringForComponent).prop("value")).toEqual("foo");
    });

    it("renders with exclusive false when multiple is true", () => {
      expect(multipleTrueWrapper.find(searchStringForComponent).prop("exclusive")).toEqual(false);
    });

    it("renders with exclusive true when multiple is false", () => {
      expect(multipleFalseWrapper.find(searchStringForComponent).prop("exclusive")).toEqual(true);
    });

    it("renders with size being large", () => {
      expect(multipleFalseWrapper.find(searchStringForComponent).prop("size")).toEqual("large");
    });

    it("renders with the default size being medium", () => {
      expect(multipleTrueWrapper.find(searchStringForComponent).prop("size")).toEqual("medium");
    });

    it("renders with the default orientation being horizontal", () => {
      expect(multipleFalseWrapper.find(searchStringForComponent).prop("orientation")).toEqual(
        "horizontal"
      );
    });

    it("renders with the orientation being vertical", () => {
      expect(multipleTrueWrapper.find(searchStringForComponent).prop("orientation")).toEqual(
        "vertical"
      );
    });
  });
});
