import React from "react";
import { Grid } from "@clutch-sh/core";
import { faComment } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";

import { LinkText } from "../cards/base";

import type { ProjectMessenger } from "./types";

const MessengerRow = ({ text, icon = faComment, url }: ProjectMessenger) => (
  <>
    <Grid item>
      <FontAwesomeIcon icon={icon} size="lg" />
    </Grid>
    <Grid item>
      <LinkText text={text} link={url} />
    </Grid>
  </>
);

export default MessengerRow;
