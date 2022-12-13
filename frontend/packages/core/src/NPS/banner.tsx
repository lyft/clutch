import React from "react";

import type { EmojiType } from "../Assets/emojis";
import { HappyEmoji } from "../Assets/emojis";
import { Button } from "../button";
import { useAppContext } from "../Contexts";
import { FeatureOn, SimpleFeatureFlag } from "../flags";
import Grid from "../grid";
import Paper from "../paper";
import styled from "../styled";
import { Typography } from "../typography";

interface BannerFeedbackProps {
  integrated?: boolean;
  icon?: EmojiType | React.ReactNode | null;
  feedbackText?: string | null;
  feedbackButtonText?: string;
  defaultOption?: string;
}

const StyledPaper = styled(Paper)({
  padding: "0px 16px",
});

/**
 * An NPS Feedback Banner which will ask the user for feedback then open up the NPSHeader
 * NOTE: requires the NPSHeader to be enabled
 *
 * @param integrated Whether this component should appear integrated into the page versus in its own container
 * @param icon The icon to display or null if none
 * @param feedbackText Feedback text, defaults to "Enjoying this feature? We would like your feedback!"
 * @param feedbackButtonText Feedback button text, defaults to "Give Feedback"
 * @returns Banner Feedback Component
 */
const BannerFeedback = ({
  integrated = false,
  icon = <HappyEmoji />,
  feedbackText = "Enjoying this feature? We would like your feedback!",
  feedbackButtonText = "Give Feedback",
  defaultOption,
}: BannerFeedbackProps) => {
  const { headerLink } = useAppContext();
  const banner = (
    <Grid container direction="row" spacing={1} alignItems="center">
      <Grid item sx={{ marginTop: "4px" }}>
        {icon}
      </Grid>
      <Grid item>
        <Typography variant="body2">{feedbackText}</Typography>
      </Grid>
      <Grid item sx={{ marginLeft: "16px" }}>
        <Button
          id="npsBannerButton"
          variant="neutral"
          text={feedbackButtonText}
          size="small"
          onClick={() => headerLink("NPS", { defaultOption })}
        />
      </Grid>
    </Grid>
  );

  return (
    <SimpleFeatureFlag feature="npsHeader">
      <FeatureOn>{integrated ? banner : <StyledPaper>{banner}</StyledPaper>}</FeatureOn>
    </SimpleFeatureFlag>
  );
};

export default BannerFeedback;
