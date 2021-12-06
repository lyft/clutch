import React from "react";
import renderer, { act } from "react-test-renderer";
import { clutch as IClutch } from "@clutch-sh/api";
import { matchers } from "@emotion/jest";
import { shallow } from "enzyme";
import { capitalize } from "lodash";

import { client } from "../../Network";
import NPSFeedback, { defaults, EmojiRatings } from "../feedback";

// Add the custom matchers provided by '@emotion/jest'
expect.extend(matchers);

describe("<NPSFeedback />", () => {
  const wizardTestResult = {
    prompt: "Test Prompt",
    freeformPrompt: "Test Freeform",
    ratingLabels: [
      {
        emoji: "SAD",
        label: "bad",
      },
      {
        emoji: "NEUTRAL",
        label: "ok",
      },
      {
        emoji: "HAPPY",
        label: "great",
      },
    ],
  };

  const clickEmoji = wrapper => {
    wrapper.find("EmojiRatings").dive().find("Styled(Component)").first().prop("onClick")(null);

    wrapper.update();
  };

  describe("basic functionality", () => {
    it("renders feedback with given origin", () => {
      const wizardWrapper = shallow(<NPSFeedback origin="WIZARD" />);
      expect(wizardWrapper).toBeTruthy();

      const anytimeWrapper = shallow(<NPSFeedback origin="ANYTIME" />);
      expect(anytimeWrapper).toBeTruthy();

      const unspecifiedWrapper = shallow(<NPSFeedback origin="ORIGIN_UNSPECIFIED" />);
      expect(unspecifiedWrapper).toBeTruthy();
    });

    describe("api success", () => {
      let wrapper;
      let useEffect;
      let spy;

      const mockUseEffect = () => {
        useEffect.mockImplementationOnce(f => f());
      };

      beforeEach(() => {
        useEffect = jest.spyOn(React, "useEffect");
        mockUseEffect();
        spy = jest.spyOn(client, "post").mockReturnValue(
          new Promise((resolve, reject) => {
            resolve({
              data: {
                originSurvey: {
                  WIZARD: wizardTestResult,
                },
              },
            });
          })
        );
        wrapper = shallow(<NPSFeedback origin="WIZARD" />);
      });

      afterEach(() => {
        wrapper.unmount();
        spy.mockClear();
      });

      it("renders survey text prompt", () => {
        expect(wrapper.find({ item: true }).at(0).find("Styled(span)").first().text()).toEqual(
          wizardTestResult.prompt
        );
      });

      it("renders emojis to <EmojiRatings />", () => {
        expect(wrapper.find("EmojiRatings").prop("ratings")).toEqual(wizardTestResult.ratingLabels);
      });

      it("renders text placeholder", () => {
        clickEmoji(wrapper);

        expect(wrapper.find("TextField").prop("placeholder")).toEqual(
          wizardTestResult.freeformPrompt
        );
      });

      it("displays a successful submission alert after submit", () => {
        wrapper.find("form").prop("onSubmit")(null);

        wrapper.update();

        const alert = wrapper.find("Styled(Alert)");

        expect(alert).toBeDefined();
        expect(alert.find("Styled(span)").text()).toEqual("Thank you for your feedback!");
      });

      it("sends feedback upon emoji selection change", () => {
        mockUseEffect();
        spy.mockClear();

        expect(spy).not.toHaveBeenCalled();
        expect(spy).toHaveBeenCalledTimes(0);

        clickEmoji(wrapper);

        expect(spy).toHaveBeenCalled();
        expect(spy).toHaveBeenCalledTimes(1);
      });
    });

    describe("api failure", () => {
      let wrapper;
      let useEffect;

      const mockUseEffect = () => {
        useEffect.mockImplementationOnce(f => f());
      };

      beforeEach(() => {
        useEffect = jest.spyOn(React, "useEffect");
        mockUseEffect();
        jest.spyOn(client, "post").mockReturnValue(
          new Promise((resolve, reject) => {
            reject(new Error("Test Error"));
          })
        );
        wrapper = shallow(<NPSFeedback origin="WIZARD" />);
      });

      afterEach(() => {
        wrapper.unmount();
      });

      it("renders default text prompt", () => {
        expect(wrapper.find({ item: true }).at(0).find("Styled(span)").first().text()).toEqual(
          defaults.prompt
        );
      });

      it("renders default emojis to <EmojiRatings />", () => {
        expect(wrapper.find("EmojiRatings").prop("ratings")).toEqual(defaults.ratingLabels);
      });

      it("renders default text placeholder", () => {
        clickEmoji(wrapper);

        expect(wrapper.find("TextField").prop("placeholder")).toEqual(defaults.freeformPrompt);
      });
    });
  });

  describe("basic rendering", () => {
    const maxLength = 180;
    let wrapper;
    let useEffect;

    const mockUseEffect = () => {
      useEffect.mockImplementationOnce(f => f());
    };

    const generateRandomString = (length, rs = "") => {
      let randomString = rs;
      randomString += Math.random().toString(20).substr(2, length);
      if (randomString.length > length) return randomString.slice(0, length);
      return generateRandomString(length, randomString);
    };

    beforeEach(() => {
      useEffect = jest.spyOn(React, "useEffect");
      mockUseEffect();
      jest.spyOn(client, "post").mockReturnValue(
        new Promise((resolve, reject) => {
          resolve({
            data: {
              originSurvey: {
                WIZARD: wizardTestResult,
              },
            },
          });
        })
      );
      wrapper = shallow(<NPSFeedback origin="WIZARD" />);
    });

    afterEach(() => {
      wrapper.unmount();
    });

    it("will not display feedback form or submit unless emoji is selected", () => {
      expect(wrapper.find({ item: true })).toHaveLength(2);
      expect(wrapper.find("Button")).toHaveLength(0);
      expect(wrapper.find("TextField")).toHaveLength(0);
    });

    it("will display text prompt at top", () => {
      expect(wrapper.find({ item: true }).at(0).find("Styled(span)")).toBeDefined();
    });

    it("will display <EmojiRatings /> below prompt", () => {
      expect(wrapper.find({ item: true }).at(1).find("EmojiRatings")).toBeDefined();
    });

    it("displays a text form and submit buttons after selection of emoji", () => {
      expect(wrapper.find({ item: true })).toHaveLength(2);

      clickEmoji(wrapper);

      expect(wrapper.find({ item: true })).toHaveLength(4);
      expect(wrapper.find("TextField")).toBeDefined();
      expect(wrapper.find("Button")).toBeDefined();
    });

    it("will update the length on feedback if input is given", () => {
      clickEmoji(wrapper);

      const testValue = "Some Feedback Text";

      let textField = wrapper.find("TextField");

      expect(textField.prop("helperText")).toEqual(`0 / ${maxLength}`);

      textField.prop("onChange")({ target: { value: testValue } });

      wrapper.update();

      textField = wrapper.find("TextField");

      expect(textField.prop("helperText")).toEqual(`${testValue.trim().length} / ${maxLength}`);
      expect(textField.prop("value")).toEqual(testValue);
    });

    it("will display an error on feedback if more input is given than maxLength", () => {
      clickEmoji(wrapper);

      const testValue = generateRandomString(200);

      let textField = wrapper.find("TextField");

      expect(textField.prop("helperText")).toEqual(`0 / ${maxLength}`);

      textField.prop("onChange")({ target: { value: testValue } });

      wrapper.update();

      textField = wrapper.find("TextField");

      expect(textField.prop("helperText")).toEqual(`${testValue.trim().length} / ${maxLength}`);
      expect(textField.prop("value")).toEqual(testValue);
      expect(textField.prop("error")).toBeTruthy();
    });

    it("will disable the submit button upon error", () => {
      clickEmoji(wrapper);

      let submitButton = wrapper.find("Button");

      expect(submitButton.prop("disabled")).toBeFalsy();

      wrapper.find("TextField").prop("onChange")({ target: { value: generateRandomString(181) } });

      wrapper.update();

      submitButton = wrapper.find("Button");

      expect(submitButton.prop("disabled")).toBeTruthy();
    });
  });

  describe("<EmojiRatings />", () => {
    const stringExample = [
      {
        emoji: "SAD",
        label: "bad",
      },
      {
        emoji: "NEUTRAL",
        label: "ok",
      },
      {
        emoji: "HAPPY",
        label: "great",
      },
    ];

    it("will display a given list of emojis and their capitalized labels", () => {
      const wrapper = shallow(<EmojiRatings ratings={stringExample} setRating={() => {}} />);
      wrapper.children().forEach((node, i) => {
        expect(node.find("Tooltip").prop("title")).toEqual(capitalize(stringExample[i].label));
        expect(node.find("Emoji").prop("type")).toEqual(stringExample[i].emoji);
      });
    });

    it("all emojis have an initial opacity of 0.5 when not selected", () => {
      const component = renderer
        .create(<EmojiRatings ratings={stringExample} setRating={() => {}} />)
        .toJSON();

      component.forEach(node => {
        expect(node).toHaveStyleRule("opacity", "0.5");
      });
    });

    it("emojis will update opacity to 1 on hover", () => {
      const component = renderer
        .create(<EmojiRatings ratings={stringExample} setRating={() => {}} />)
        .toJSON();

      component.forEach(node => {
        expect(node).toHaveStyleRule("opacity", "1", { target: ":hover" });
      });
    });

    it("emojis wll update opacity to 1 on selection", () => {
      let component;

      act(() => {
        component = renderer.create(<EmojiRatings ratings={stringExample} setRating={() => {}} />);
      });

      let [firstEmoji] = component.toJSON();

      expect(firstEmoji).toHaveStyleRule("opacity", "0.5");
      expect(firstEmoji.props.selected).toBeFalsy();

      act(() => {
        firstEmoji.props.onClick();
      });

      [firstEmoji] = component.toJSON();

      expect(firstEmoji).toHaveStyleRule("opacity", "1");
      expect(firstEmoji.props.selected).toBeTruthy();
    });

    it("will fetch emojis based on integers with a given enum", () => {
      const enumExample = [
        { emoji: 1, label: "bad" },
        { emoji: 2, label: "ok" },
        { emoji: 3, label: "great" },
      ];

      const enums = IClutch.feedback.v1.EmojiRating;
      const wrapper = shallow(<EmojiRatings ratings={enumExample} setRating={() => {}} />);

      wrapper.children().forEach((node, i) => {
        expect(enums[node.find("Emoji").prop("type")]).toEqual(enumExample[i].emoji);
      });
    });

    it("will return a given emoji on select", () => {
      let selected = null;

      const wrapper = shallow(
        <EmojiRatings
          ratings={stringExample}
          setRating={rating => {
            selected = rating;
          }}
        />
      );

      expect(selected).toBeNull();

      const neutralEmoji = wrapper.find("Tooltip").at(1).find("Styled(Component)");

      neutralEmoji.prop("onClick")(null);

      wrapper.update();

      expect(selected).toEqual(stringExample[1]);
    });
  });

  // Verifies layout changes for given origins
  describe("Wizard Origin Rendering", () => {
    let wrapper;
    let useEffect;

    const mockUseEffect = () => {
      useEffect.mockImplementationOnce(f => f());
    };

    beforeEach(() => {
      useEffect = jest.spyOn(React, "useEffect");
      mockUseEffect();
      jest.spyOn(client, "post").mockReturnValue(
        new Promise((resolve, reject) => {
          resolve({
            data: {
              originSurvey: {
                WIZARD: wizardTestResult,
              },
            },
          });
        })
      );
      wrapper = shallow(<NPSFeedback origin="WIZARD" />);
    });

    afterEach(() => {
      wrapper.unmount();
    });

    it("center aligns the grid", () => {
      expect(wrapper.find("form").childAt(0).prop("justify")).toEqual("center");
    });

    it("styles the submit button correctly", () => {
      wrapper.find("EmojiRatings").dive().find("Styled(Component)").first().prop("onClick")(null);

      wrapper.update();

      expect(wrapper.find("Button").prop("variant")).toEqual("secondary");
    });
  });
});
