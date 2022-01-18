import React, { useState } from "react";
import { clutch as IClutch } from "@clutch-sh/api";
import { capitalize, isInteger } from "lodash";

import Emoji, { EmojiType } from "../Assets/emojis";
import type { IconButtonSize } from "../button";
import { IconButton } from "../button";
import { Tooltip } from "../Feedback";
import styled from "../styled";

export type Rating = {
  emoji: string;
  label: string;
};

type EmojiRatingsProps = {
  ratings: IClutch.feedback.v1.IRatingLabel[];
  setRating: (Rating) => void;
  placement?: "top" | "bottom";
  size?: IconButtonSize;
};

// Will convert a given integer to a typed enum key if necessary
const getKey = (map, val): string => Object.keys(map).find(key => map[key] === val);

/**
 * EmojiRatings component which will take an array of emojis and given ratings and create IconButtons with them and return them on selection
 *
 * @param ratings given array of ratings to display
 * @param setRating function which will return the given selected rating
 * @returns rendered EmojiRatings component
 */
const EmojiRatings = ({
  ratings = [],
  setRating,
  placement = "top",
  size = "small",
}: EmojiRatingsProps) => {
  const [selectedRating, selectRating] = useState<Rating>(null);

  const StyledIconButton = styled(IconButton)<{
    $selected: boolean;
  }>(
    {
      "&:hover": {
        opacity: 1,
      },
    },
    props => ({
      opacity: props.$selected ? 1 : 0.5,
    })
  );

  const select = (rating: Rating) => {
    selectRating(rating);
    setRating(rating);
  };

  return (
    <>
      {ratings.map((rating: IClutch.feedback.v1.IRatingLabel) => {
        const { label } = rating;
        const emoji = isInteger(rating.emoji)
          ? getKey(IClutch.feedback.v1.EmojiRating, rating.emoji)
          : rating.emoji;

        return (
          <Tooltip key={label} title={capitalize(label)} placement={placement}>
            <StyledIconButton
              key={`rating-${emoji}`}
              variant="neutral"
              $selected={selectedRating?.label === label}
              onClick={() => select({ label, emoji: emoji as string })}
              size={size}
            >
              <Emoji type={emoji as EmojiType} size={size} />
            </StyledIconButton>
          </Tooltip>
        );
      })}
    </>
  );
};

export default EmojiRatings;
