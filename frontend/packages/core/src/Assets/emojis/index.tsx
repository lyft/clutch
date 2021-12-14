import React from "react";

import type { SVGProps } from "../global";

import HappyEmoji from "./happy";
import NeutralEmoji from "./neutral";
import SadEmoji from "./sad";

const EmojiTypes = {
  HAPPY: ({ ...props }) => <HappyEmoji {...props} />,
  NEUTRAL: ({ ...props }) => <NeutralEmoji {...props} />,
  SAD: ({ ...props }) => <SadEmoji {...props} />,
};

export type EmojiType = keyof typeof EmojiTypes;

interface EmojiProps extends SVGProps {
  type: EmojiType;
}

/**
 * Shorthand component to return an Emoji based on a given input type
 *
 * @param type the given type of emoji to return from EmojiType
 * @returns the given emoji extended with any additional properties passed in
 */
const Emoji = ({ type, ...props }: EmojiProps) => EmojiTypes[type](props);

export default Emoji;

export { HappyEmoji, NeutralEmoji, SadEmoji };
