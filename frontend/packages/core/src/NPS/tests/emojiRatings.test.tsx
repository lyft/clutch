import React from "react";
import renderer, { act } from "react-test-renderer";
import { clutch as IClutch } from "@clutch-sh/api";
import { matchers } from "@emotion/jest";
import { shallow } from "enzyme";
import { capitalize } from "lodash";

import EmojiRatings from "../emojiRatings";

// Add the custom matchers provided by '@emotion/jest'
expect.extend(matchers);

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
