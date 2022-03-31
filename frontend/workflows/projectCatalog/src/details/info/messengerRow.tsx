import React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";
import { Grid } from "@clutch-sh/core";
import type { IconProp } from "@fortawesome/fontawesome-svg-core";
import { faSlack } from "@fortawesome/free-brands-svg-icons";
import { faComment } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";

import { LinkText } from "../helpers";

const MessengerRow = ({ projectData }: { projectData: IClutch.core.project.v1.IProject }) => {
  const [text, setText] = React.useState<string>();
  const [link, setLink] = React.useState<string>();
  const [icon, setIcon] = React.useState<IconProp>(faComment);

  React.useEffect(() => {
    if (projectData) {
      // TODO (jslaughter): Should be customizable
      const slackLink = (projectData?.linkGroups || []).find(lg => lg.name === "Slack");

      if (projectData?.data?.slack) {
        setText(projectData?.data?.slack as string);
        setIcon(faSlack);
      }

      if (slackLink && slackLink?.links?.length) {
        setLink(slackLink?.links[0].url as string);
      }
    }
  }, [projectData]);

  return (
    <Grid container item spacing={1}>
      <Grid item>
        <FontAwesomeIcon icon={icon} size="lg" />
      </Grid>
      <Grid item>{text && <LinkText text={text} link={link} />}</Grid>
    </Grid>
  );
};

export default MessengerRow;
