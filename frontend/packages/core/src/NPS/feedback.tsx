import React, { useState } from "react";
import { clutch as IClutch } from "@clutch-sh/api";
import styled from "@emotion/styled";
import { Grid as MuiGrid } from "@material-ui/core";
import MuiSuccessIcon from "@material-ui/icons/CheckCircle";
import { debounce } from "lodash";
import { v4 as uuid } from "uuid";

import { userId } from "../AppLayout/user";
import { Button } from "../button";
import { Alert } from "../Feedback";
import { TextField } from "../Input";
import { client } from "../Network";
import type { ClutchError } from "../Network/errors";

import EmojiRatings, { Rating } from "./emojiRatings";

/** Interfaces */

interface FeedbackOptions {
  origin: "ORIGIN_UNSPECIFIED" | "WIZARD" | "ANYTIME";
  onSubmit?: (submit: boolean) => void;
}

// Defaults in case of API failure
export const defaults: IClutch.feedback.v1.ISurvey = {
  prompt: "Rate Your Experience using Clutch",
  freeformPrompt: "What would you recommend to improve this?",
  ratingLabels: [
    {
      emoji: 1,
      label: "bad",
    },
    {
      emoji: 2,
      label: "ok",
    },
    {
      emoji: 3,
      label: "great",
    },
  ],
};

/** Styling */

const Text = styled.span({
  fontWeight: "bold",
  color: "#0D1030",
  fontSize: "14px",
});

const StyledGrid = styled(MuiGrid)({
  borderRadius: "8px",
  padding: "20px",
});

const StyledAlert = styled(Alert)({
  margin: "32px",
  alignItems: "center",
});

/**
 * NPS feedback component which is the base for both Wizard and Anytime
 * Will fetch given survey options from the server based on the given origin
 * Then display a feedback component based on the given emoji ratings
 *
 * @param opts Available feedback options
 * @returns NPSFeedback component
 */
const NPSFeedback = (opts: FeedbackOptions = { origin: "ORIGIN_UNSPECIFIED" }) => {
  const [hasSubmit, setSubmit] = useState<boolean>(false);
  const [selected, setSelected] = useState<Rating>(null);
  const [feedback, setFeedback] = useState<string>("");
  const [error, setError] = useState<boolean>(false);
  const [survey, setSurvey] = useState<IClutch.feedback.v1.ISurvey>({});
  const [requestId, setRequestId] = useState<string>("");
  const maxLength = 180;
  const debounceTimer = 500;

  const trimmed =
    feedback.trim().length > maxLength
      ? `${feedback.trim().substring(0, maxLength - 3)}...`
      : feedback;

  /** Property objects used to extend components and remove console warnings */
  const AlertProps = {
    iconMapping: {
      info: <MuiSuccessIcon style={{ color: "#3548d4" }} />,
    },
  };

  const textFieldProps = {
    fullWidth: true,
    InputProps: {
      rows: 3,
      rowsMax: 3,
    },
    style: {
      marginTop: "15px",
    },
  };

  // Will fetch the survey results for the given origin on load
  React.useEffect(() => {
    let data: IClutch.feedback.v1.ISurvey = defaults;

    client
      .post("/v1/feedback/getSurveys", {
        origins: [opts.origin],
      })
      .then(response => {
        const surveyData: IClutch.feedback.v1.IGetSurveysResponse = response?.data?.originSurvey;

        data = surveyData[opts.origin] ?? defaults;
      })
      .catch((err: ClutchError) => {
        // eslint-disable-next-line no-console
        console.error(err);
      })
      .finally(() => {
        setRequestId(uuid());
        setSurvey(data);
      });
  }, []);

  // Will debounce feedback requests to the server in case of multiple quick changes to selected
  const sendFeedback = React.useCallback(
    debounce((formData: IClutch.feedback.v1.ISubmitFeedbackRequest) => {
      client
        .post("/v1/feedback/submitFeedback", { userId: userId(), ...formData })
        .catch((err: ClutchError) => {
          // eslint-disable-next-line no-console
          console.error(err);
        });
    }, debounceTimer),
    []
  );

  // On a change to submit or selected will attempt to send a feedback request
  React.useEffect(() => {
    if (selected) {
      sendFeedback({
        id: requestId,
        feedback: {
          ratingLabel: selected.label,
          ratingScale: {
            emoji: IClutch.feedback.v1.EmojiRating[selected.emoji],
          },
          urlPath: window.location.pathname,
          freeformResponse: trimmed,
        },
        metadata: {
          origin: IClutch.feedback.v1.Origin[opts.origin],
          userSubmitted: hasSubmit,
          survey,
        },
      });
    }
  }, [selected, hasSubmit]);

  // Form onSubmit handler
  const submitFeedback = e => {
    if (e) {
      e.preventDefault();
    }
    setSubmit(true);
    if (opts.onSubmit) {
      opts.onSubmit(true);
    }
  };

  return (
    <>
      {hasSubmit ? (
        <StyledAlert severity="info" {...AlertProps}>
          <Text>Thank you for your feedback!</Text>
        </StyledAlert>
      ) : (
        <form onSubmit={submitFeedback}>
          <StyledGrid container direction="row" alignItems="center" justify="center">
            <MuiGrid item xs={6}>
              <Text>{survey.prompt}</Text>
            </MuiGrid>
            <MuiGrid item xs={6} style={{ display: "flex", justifyContent: "space-around" }}>
              <EmojiRatings ratings={survey.ratingLabels} setRating={setSelected} />
            </MuiGrid>
            {selected !== null && (
              <>
                <MuiGrid item xs={12}>
                  <TextField
                    multiline
                    placeholder={survey.freeformPrompt}
                    value={feedback}
                    helperText={`${feedback?.trim().length} / ${maxLength}`}
                    error={error}
                    onChange={e => {
                      setFeedback(e?.target?.value);
                      setError(e?.target?.value?.trim().length > maxLength);
                    }}
                    {...textFieldProps}
                  />
                </MuiGrid>
                <MuiGrid item xs={4}>
                  <Button type="submit" text="Submit" variant="secondary" disabled={error} />
                </MuiGrid>
              </>
            )}
          </StyledGrid>
        </form>
      )}
    </>
  );
};

export default NPSFeedback;
