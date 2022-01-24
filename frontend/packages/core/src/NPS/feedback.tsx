import React, { useState } from "react";
import { clutch as IClutch } from "@clutch-sh/api";
import { Grid as MuiGrid } from "@material-ui/core";
import MuiSuccessIcon from "@material-ui/icons/CheckCircle";
import { debounce } from "lodash";
import { v4 as uuid } from "uuid";

import { userId } from "../AppLayout/user";
import { Button } from "../button";
import { Alert } from "../Feedback";
import type { SelectOption } from "../Input";
import { Select, TextField } from "../Input";
import { client } from "../Network";
import type { ClutchError } from "../Network/errors";
import styled from "../styled";
import { Typography } from "../typography";

import EmojiRatings, { Rating } from "./emojiRatings";

/** Interfaces */

type Origins = "WIZARD" | "HEADER";

interface FeedbackOptions {
  origin: Origins;
  feedbackTypes?: SelectOption[];
  onSubmit?: (submit: boolean) => void;
}

// Defaults in case of API failure
export const defaults: IClutch.feedback.v1.ISurvey = {
  prompt: "Rate Your Experience",
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

const StyledButton = styled(Button)<{ $origin: Origins }>({}, props =>
  props.$origin === "WIZARD"
    ? {
        fontSize: "14px",
        padding: "0 8px",
        height: "18px",
      }
    : null
);

const StyledTextField = styled(TextField)<{ $origin: Origins }>(
  {
    margin: "16px 0px 32px 0px",
  },
  props => ({
    ".MuiInputBase-root": {
      fontSize: props.$origin === "WIZARD" ? "14px" : "16px",
    },
  })
);

const FeedbackAlert = () => {
  const AlertProps = {
    iconMapping: {
      info: <MuiSuccessIcon style={{ color: "#3548d4" }} />,
    },
    style: {
      margin: "32px",
      alignItems: "center",
    },
  };

  return (
    <Alert severity="info" {...AlertProps}>
      <Typography variant="subtitle3">Thank you for your feedback!</Typography>
    </Alert>
  );
};

/**
 * NPS feedback component which is the base for both Wizard and Anytime.
 * Will fetch given survey options from the server based on the given origin
 * Then display a feedback component based on the given emoji ratings
 *
 * @param opts Available feedback options
 * @returns NPSFeedback component
 */
const NPSFeedback = ({ origin = "HEADER", ...options }: FeedbackOptions) => {
  const [hasSubmit, setHasSubmit] = useState<boolean>(false);
  const [selectedEmoji, setSelectedEmoji] = useState<Rating>(null);
  const [freeformFeedback, setFreeformFeedback] = useState<string>("");
  const [error, setError] = useState<boolean>(false);
  const [survey, setSurvey] = useState<IClutch.feedback.v1.ISurvey>({});
  const [feedbackType, setFeedbackType] = useState<string>(null);
  const [requestId, setRequestId] = useState<string>("");
  const maxLength = 180;
  const debounceTimer = 500;
  const wizardOrigin = origin === "WIZARD";

  const trimmed =
    freeformFeedback.trim().length > maxLength
      ? `${freeformFeedback.trim().substring(0, maxLength - 3)}...`
      : freeformFeedback;

  const textFieldProps = {
    fullWidth: true,
    InputProps: {
      rows: 3,
      rowsMax: 3,
    },
  };

  // Will fetch the survey results for the given origin on load
  React.useEffect(() => {
    let data: IClutch.feedback.v1.ISurvey = defaults;

    client
      .post("/v1/feedback/getSurveys", {
        origins: [origin],
      })
      .then(response => {
        const surveyData: IClutch.feedback.v1.IGetSurveysResponse = response?.data?.originSurvey;

        data = surveyData[origin] ?? defaults;
      })
      .catch((err: ClutchError) => {
        // eslint-disable-next-line no-console
        console.error(err);
      })
      .finally(() => {
        setRequestId(uuid());
        setSurvey(data);

        if (options.feedbackTypes && options.feedbackTypes.length) {
          setFeedbackType(options.feedbackTypes[0].value || options.feedbackTypes[0].label);
        }
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
    if (selectedEmoji) {
      sendFeedback({
        id: requestId,
        feedback: {
          feedbackType,
          freeformResponse: trimmed,
          ratingLabel: selectedEmoji.label,
          ratingScale: {
            emoji: IClutch.feedback.v1.EmojiRating[selectedEmoji.emoji],
          },
          urlPath: window.location.pathname,
        },
        metadata: {
          survey,
          origin: IClutch.feedback.v1.Origin[origin],
          userSubmitted: hasSubmit,
          urlSearchParams: window.location.search,
        },
      });
    }
  }, [selectedEmoji, feedbackType, hasSubmit]);

  // Form onSubmit handler
  const submitFeedback = e => {
    if (e) {
      e.preventDefault();
    }
    setHasSubmit(true);
    if (options.onSubmit) {
      options.onSubmit(true);
    }
  };

  if (hasSubmit) {
    return <FeedbackAlert />;
  }

  return (
    <form onSubmit={submitFeedback}>
      <MuiGrid
        container
        direction="row"
        alignItems="center"
        style={{ padding: wizardOrigin ? "16px" : "24px" }}
      >
        <MuiGrid item xs>
          <Typography variant={wizardOrigin ? "subtitle3" : "subtitle2"}>
            {survey.prompt}
          </Typography>
        </MuiGrid>
        <MuiGrid
          item
          xs={wizardOrigin ? 6 : 12}
          style={{ display: "flex", justifyContent: "space-around" }}
        >
          <EmojiRatings
            ratings={survey.ratingLabels}
            setRating={setSelectedEmoji}
            placement={wizardOrigin ? "top" : "bottom"}
            buttonSize={wizardOrigin ? "small" : "medium"}
          />
        </MuiGrid>
        {selectedEmoji !== null && (
          <>
            {!wizardOrigin && options.feedbackTypes && (
              <MuiGrid item xs={12} style={{ margin: "32px 0px 16px 0px" }}>
                <Select
                  name="anytimeSelect"
                  label="Choose a type of feedback you want to submit"
                  options={options.feedbackTypes}
                  onChange={setFeedbackType}
                />
              </MuiGrid>
            )}
            <MuiGrid item xs={12}>
              <StyledTextField
                multiline
                fullWidth
                $origin={origin}
                placeholder={survey.freeformPrompt}
                value={freeformFeedback}
                helperText={`${freeformFeedback?.trim().length} / ${maxLength}`}
                error={error}
                onChange={e => {
                  setFreeformFeedback(e?.target?.value);
                  setError(e?.target?.value?.trim().length > maxLength);
                }}
                {...textFieldProps}
              />
            </MuiGrid>
            <MuiGrid
              item
              xs={12}
              style={{
                display: "flex",
                justifyContent: wizardOrigin ? "center" : "flex-end",
              }}
            >
              <StyledButton
                $origin={origin}
                type="submit"
                text="Submit"
                variant={wizardOrigin ? "secondary" : "primary"}
                disabled={error}
              />
            </MuiGrid>
          </>
        )}
      </MuiGrid>
    </form>
  );
};

export default NPSFeedback;
