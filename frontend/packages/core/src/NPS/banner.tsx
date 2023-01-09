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

export interface FeedbackBannerProps {
  /**
   * Whether this component should appear integrated into the page versus in its own container
   * @defaultValue false
   */
  integrated?: boolean;
  /**
   * The icon to display or null if none
   * @defaultValue <HappyEmoji />
   */
  icon?: EmojiType | React.ReactNode | null;
  /**
   * Feedback text
   * @defaultValue "Enjoying this feature? We would like your feedback!"
   */
  feedbackText?: string | null;
  /**
   * Feedback Button Text
   * @defaultValue "Give Feedback"
   */
  feedbackButtonText?: string;
  /**
   * Default feedback option in dropdown to select
   * @defaultValue ""
   */
  defaultFeedbackOption?: string;
}

const StyledPaper = styled(Paper)({
  padding: "0px 16px",
  "@media screen and (max-width: 970px)": {
    padding: "8px 16px",
  },
});

/**
 * An NPS Feedback Banner which will open up the NPSHeader component to ask for feedback
 */
export const Banner = ({
  integrated = false,
  icon = <HappyEmoji />,
  feedbackText = "Enjoying this feature? We would like your feedback!",
  feedbackButtonText = "Give Feedback",
  defaultFeedbackOption,
}: FeedbackBannerProps) => {
  const { triggerHeaderItem } = useAppContext();
  const banner = (
    <Grid container direction="row" spacing={1} alignItems="center" justifyContent="center">
      {icon && (
        <Grid item>
          <div data-testid="nps-banner-icon" style={{ display: "flex" }}>
            {icon}
          </div>
        </Grid>
      )}
      <Grid item>
        <Typography data-testid="nps-banner-text" variant="body2">
          {feedbackText}
        </Typography>
      </Grid>
      <Grid item sx={{ marginLeft: "16px" }}>
        <Button
          id="nps-banner-button"
          data-testid="nps-banner-button"
          variant="neutral"
          text={feedbackButtonText}
          size="small"
          onClick={() => triggerHeaderItem("NPS", { defaultFeedbackOption })}
        />
      </Grid>
    </Grid>
  );

  return integrated ? (
    banner
  ) : (
    <StyledPaper data-testid="nps-banner-container">{banner}</StyledPaper>
  );
};

/**
 * An NPS Feedback Banner which will ask the user for feedback then open up the NPSHeader
 * NOTE: requires the NPSHeader to be enabled
 */
const FeedbackBanner = ({ ...props }: FeedbackBannerProps) => (
  <SimpleFeatureFlag feature="npsHeader">
    <FeatureOn>
      <Banner {...props} />
    </FeatureOn>
  </SimpleFeatureFlag>
);

export default FeedbackBanner;
