import React from "react";
import { matchers } from "@emotion/jest";
import { shallow } from "enzyme";

import contextValues from "../../Contexts/tests/testContext";
import { client } from "../../Network";
import NPSFeedback, { defaults, FEEDBACK_MAX_LENGTH } from "../feedback";
import { generateFeedbackTypes } from "../header";

// Adds the custom matchers provided by '@emotion/jest'
expect.extend(matchers);

const feedbackTypes = generateFeedbackTypes(contextValues.workflows);

describe("<NPSFeedback />", () => {
  const defaultResult = {
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
                  WIZARD: defaultResult,
                  HEADER: defaultResult,
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
        expect(wrapper.find({ item: true }).at(0).find("Typography").childAt(0).text()).toEqual(
          defaultResult.prompt
        );
      });
      it("renders emojis to <EmojiRatings />", () => {
        expect(wrapper.find("EmojiRatings").prop("ratings")).toEqual(defaultResult.ratingLabels);
      });
      it("renders text placeholder", () => {
        clickEmoji(wrapper);
        expect(wrapper.find("Styled(TextField)").prop("placeholder")).toEqual(
          defaultResult.freeformPrompt
        );
      });
      it("displays a successful submission alert after submit", () => {
        wrapper.find("form").prop("onSubmit")(null);
        wrapper.update();
        const alert = wrapper.find("FeedbackAlert");
        expect(alert).toBeDefined();
        expect(alert.dive().find("Typography").childAt(0).text()).toBe(
          "Thank you for your feedback!"
        );
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
        expect(wrapper.find({ item: true }).at(0).find("Typography").childAt(0).text()).toEqual(
          defaults.prompt
        );
      });
      it("renders default emojis to <EmojiRatings />", () => {
        expect(wrapper.find("EmojiRatings").prop("ratings")).toEqual(defaults.ratingLabels);
      });
      it("renders default text placeholder", () => {
        clickEmoji(wrapper);
        expect(wrapper.find("Styled(TextField)").prop("placeholder")).toEqual(
          defaults.freeformPrompt
        );
      });
    });
  });

  describe("basic rendering", () => {
    const maxLength = FEEDBACK_MAX_LENGTH;
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
                WIZARD: defaultResult,
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
      expect(wrapper.find("Styled(TextField)")).toHaveLength(0);
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
      expect(wrapper.find("Styled(TextField)")).toBeDefined();
      expect(wrapper.find("Styled(Button)")).toBeDefined();
    });

    it("will update the length on feedback if input is given", () => {
      clickEmoji(wrapper);

      const testValue = "Some Feedback Text";

      let textField = wrapper.find("Styled(TextField)");

      expect(textField.prop("helperText")).toBe(`0 / ${maxLength}`);

      textField.prop("onChange")({ target: { value: testValue } });

      wrapper.update();

      textField = wrapper.find("Styled(TextField)");

      expect(textField.prop("helperText")).toBe(`${testValue.trim().length} / ${maxLength}`);
      expect(textField.prop("value")).toEqual(testValue);
    });

    it("will display an error on feedback if more input is given than maxLength", () => {
      clickEmoji(wrapper);

      const testValue = generateRandomString(FEEDBACK_MAX_LENGTH * 2);

      let textField = wrapper.find("Styled(TextField)");

      expect(textField.prop("helperText")).toBe(`0 / ${maxLength}`);

      textField.prop("onChange")({ target: { value: testValue } });

      wrapper.update();

      textField = wrapper.find("Styled(TextField)");

      expect(textField.prop("helperText")).toBe(`${testValue.trim().length} / ${maxLength}`);
      expect(textField.prop("value")).toEqual(testValue);
      expect(textField.prop("error")).toBeTruthy();
    });

    it("will disable the submit button upon error", () => {
      clickEmoji(wrapper);

      let submitButton = wrapper.find("Styled(Button)");

      expect(submitButton.prop("disabled")).toBeFalsy();

      wrapper.find("Styled(TextField)").prop("onChange")({
        target: { value: generateRandomString(FEEDBACK_MAX_LENGTH + 1) },
      });

      wrapper.update();

      submitButton = wrapper.find("Styled(Button)");

      expect(submitButton.prop("disabled")).toBeTruthy();
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
                WIZARD: defaultResult,
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

    it("renders", () => {
      expect(wrapper).toBeDefined();
    });

    it("styles the submit button correctly", () => {
      wrapper.find("EmojiRatings").dive().find("Styled(Component)").first().prop("onClick")(null);

      wrapper.update();

      expect(wrapper.find("Styled(Button)").prop("variant")).toBe("secondary");
    });
  });

  describe("Header Origin Rendering", () => {
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
                HEADER: defaultResult,
              },
            },
          });
        })
      );
      wrapper = shallow(<NPSFeedback origin="HEADER" feedbackTypes={feedbackTypes} />);
    });

    afterEach(() => {
      wrapper.unmount();
    });

    it("renders", () => {
      expect(wrapper).toBeDefined();
    });

    it("styles the submit button correctly", () => {
      clickEmoji(wrapper);

      expect(wrapper.find("Styled(Button)").prop("variant")).toBe("primary");
    });
  });
});
