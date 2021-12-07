import React, { useState } from "react";
import { clutch as IClutch } from "@clutch-sh/api";
import styled from "@emotion/styled";
import { capitalize, isInteger } from "lodash";

import Emoji, { EmojiType } from "../Assets/emojis";
import { IconButton } from "../button";
import { Tooltip } from "../Feedback";

export type Rating = {
  emoji: string;
  label: string;
};

/**
 * EmojiRatings component which will take an array of emojis and given ratings and create IconButtons with them and return them on selection
 *
 * @param ratings given array of ratings to display
 * @param setRating function which will return the given selected rating
 * @param size allows overriding of a given size, "large" is the only appropriate value
 * @returns rendered EmojiRatings component
 */
const EmojiRatings = ({ ratings = [], setRating }) => {
  const [selectedRating, selectRating] = useState<Rating>(null);

  const StyledIconButton = styled(IconButton)<{
    selected: boolean;
  }>(
    {
      "&:hover": {
        opacity: 1,
      },
    },
    props => ({
      opacity: props.selected ? 1 : 0.5,
    })
  );

  const select = (rating: Rating) => {
    selectRating(rating);
    setRating(rating);
  };

  // Will convert a given integer to a typed enum key if necessary
  const getKey = (map, val) => Object.keys(map).find(key => map[key] === val);

  return (
    <>
      {ratings.map((rating: Rating) => {
        const { label } = rating;
        let { emoji } = rating;

        if (isInteger(emoji)) {
          emoji = getKey(IClutch.feedback.v1.EmojiRating, emoji);
        }

        return (
          <Tooltip key={label} title={capitalize(label)} placement="top">
            <StyledIconButton
              key={`rating-${emoji}`}
              variant="neutral"
              size="small"
              selected={selectedRating?.label === label}
              onClick={() => select(rating)}
            >
              <Emoji type={emoji as EmojiType} />
            </StyledIconButton>
          </Tooltip>
        );
      })}
    </>
  );
};

export default EmojiRatings;
