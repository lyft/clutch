import React from "react";

import type { SVGProps } from "../global";

import HappyEmoji from "./happy";
import NeutralEmoji from "./neutral";
import SadEmoji from "./sad";

export type EmojiType = "HAPPY" | "NEUTRAL" | "SAD";

interface EmojiProps extends SVGProps {
  type: EmojiType;
}

const Emoji = ({ type, ...props }: EmojiProps) => {
  switch (type) {
    case "HAPPY":
      return <HappyEmoji {...props} />;
    case "NEUTRAL":
      return <NeutralEmoji {...props} />;
    case "SAD":
      return <SadEmoji {...props} />;
    default:
      return null;
  }
};

export default Emoji;

export { HappyEmoji, NeutralEmoji, SadEmoji };
