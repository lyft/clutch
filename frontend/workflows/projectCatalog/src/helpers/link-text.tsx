import React from "react";
import { Link, Typography } from "@clutch-sh/core";

const LinkText = ({ text, link }: { text: string; link?: string }) => {
  const returnText = <Typography variant="body2">{text}</Typography>;

  if (link && text) {
    return (
      <Link href={link} whiteSpace="nowrap">
        {returnText}
      </Link>
    );
  }

  return returnText;
};

export default LinkText;
