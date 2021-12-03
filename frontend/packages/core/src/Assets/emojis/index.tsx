import React from "react";

import HappyEmoji from "./happy";
import NeutralEmoji from "./neutral";
import SadEmoji from "./sad";

const Emoji = ({ type, ...props }) => {
  switch (type ?? type.toUpperCase()) {
    case "HAPPY":
      return <HappyEmoji {...props} />;
    case "NEUTRAL":
      return <NeutralEmoji {...props} />;
    case "SAD":
      return <SadEmoji {...props} />;
    default:
      throw new Error(`Emoji '${type}' is an invalid type`);
  }
};

export default Emoji;

export { HappyEmoji, NeutralEmoji, SadEmoji };
