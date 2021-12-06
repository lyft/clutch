import React, { useCallback, useEffect, useState } from "react";
import { clutch as IClutch } from "@clutch-sh/api";
import styled from "@emotion/styled";
import { Grid as MuiGrid } from "@material-ui/core";
import MuiSuccessIcon from "@material-ui/icons/CheckCircle";
import { capitalize, debounce, isInteger } from "lodash";
import { v4 as uuid } from "uuid";

import { userId } from "../AppLayout/user";
import Emoji from "../Assets/emojis";
import { Button, IconButton } from "../button";
import { Alert, Tooltip } from "../Feedback";
import { Select, TextField } from "../Input";
import { client } from "../Network";
import type { ClutchError } from "../Network/errors";

/** Interfaces */

interface FeedbackOptions {
  origin: "ORIGIN_UNSPECIFIED" | "WIZARD" | "ANYTIME";
  onSubmit?: (submit: boolean) => void;
}

type Rating = {
  emoji: string;
  label: string;
};

// TODO: (jslaughter) update with milestone 2 anytime typing
interface Survey extends IClutch.feedback.v1.ISurvey {
  freeformLabel?: string;
  feedbackTypeLabel?: string;
  feedbackType?: { label: string }[];
}

// Defaults in case of API failure
const defaults: Survey = {
  prompt: "Rate Your Experience using Clutch",
  freeformPrompt: "What would you recommend to improve this workflow?",
  freeformLabel: "Do you have any thought's you'd like to share? (optional)",
  feedbackTypeLabel: "Choose a type of feedback you want to submit",
  feedbackType: [{ label: "General" }, { label: "Dash" }],
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

const Text = styled.span<{ origin: string }>(
  {
    fontWeight: "bold",
    color: "#0D1030",
  },
  props => ({
    fontSize: props.origin === "ANYTIME" ? "16px" : "14px",
  })
);

const StyledGrid = styled(MuiGrid)({
  borderRadius: "8px",
  padding: "20px",
});

const StyledAlert = styled(Alert)({
  margin: "32px",
  alignItems: "center",
});

/** Components */

/**
 * EmojiRatings component which will take an array of emojis and given ratings and create IconButtons with them and return them on selection
 *
 * @param ratings given array of ratings to display
 * @param setRating function which will return the given selected rating
 * @param size allows overriding of a given size, "large" is the only appropriate value
 * @returns rendered EmojiRatings component
 */
const EmojiRatings = ({ ratings = [], setRating, size }) => {
  const [selectedRating, selectRating] = useState<Rating>(null);

  // TODO: (jslaughter) Update with large sizing with material-ui@5
  const StyledIconButton = styled(IconButton)<{
    selected: boolean;
    overridesize: string;
  }>(
    {
      "&:hover": {
        opacity: 1,
      },
    },
    props => ({
      height: props.overridesize === "large" ? "60px" : "48px",
      width: props.overridesize === "large" ? "60px" : "48px",
      opacity: props.selected ? 1 : 0.5,
    })
  );

  const select = (rating: Rating) => {
    selectRating(rating);
    setRating(rating);
  };

  // Will convert a given integer to a typed enum key
  const getKey = (map, val) => Object.keys(map).find(key => map[key] === val);

  return (
    <>
      {ratings.map((rating: Rating) => {
        const { label } = rating;
        let { emoji } = rating;

        if (isInteger(emoji)) {
          emoji = getKey(IClutch.feedback.v1.EmojiRating, emoji);
        }

        return (
          <Tooltip key={label} title={capitalize(label)}>
            <StyledIconButton
              key={`rating-${emoji}`}
              variant="neutral"
              size="medium"
              selected={selectedRating?.label === label}
              overridesize={size}
              onClick={() => select(rating)}
            >
              <Emoji type={emoji} size="large" />
            </StyledIconButton>
          </Tooltip>
        );
      })}
    </>
  );
};

const NPSFeedback = (opts: FeedbackOptions = { origin: "ORIGIN_UNSPECIFIED" }) => {
  const [hasSubmit, setSubmit] = useState<boolean>(false);
  const [selected, setSelected] = useState<Rating>(null);
  const [feedback, setFeedback] = useState<string>("");
  const [error, setError] = useState<boolean>(false);
  // const [survey, setSurvey] = useState<IClutch.feedback.v1.ISurvey>({});
  const [type, setType] = useState<string>(null);
  const [survey, setSurvey] = useState<Survey>({});
  const requestId = uuid();
  const maxLength = 180;
  const debounceTimer = 500;

  const trimmed =
    feedback.trim().length > maxLength
      ? `${feedback.trim().substring(0, maxLength - 3)}...`
      : feedback;

  /** Property objects used to extend components and remove unnecessary console warnings */
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
  useEffect(() => {
    // let data: IClutch.feedback.v1.ISurvey = defaults;
    let data: Survey = defaults;

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
        setSurvey(data);
      });
  }, []);

  // Will debounce feedback requests to the server in case of multiple quick changes to selected
  const sendFeedback = useCallback(
    debounce((formData: IClutch.feedback.v1.ISubmitFeedbackRequest) => {
      client.post("/v1/feedback/submitFeedback", { userId: userId(), ...formData }).catch(e => {
        // eslint-disable-next-line no-console
        console.error(e);
      });
    }, debounceTimer),
    []
  );

  // On a change to submit or selected will attempt to send a feedback request
  useEffect(() => {
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
          feedbackType: type,
        },
        metadata: {
          origin: IClutch.feedback.v1.Origin[opts.origin],
          userSubmitted: hasSubmit,
          survey,
        },
      });
    }
  }, [selected, hasSubmit]);

  const submitFeedback = e => {
    e.preventDefault();
    setSubmit(true);
    if (opts.onSubmit) {
      opts.onSubmit(true);
    }
  };

  return (
    <>
      {hasSubmit ? (
        <StyledAlert severity="info" {...AlertProps}>
          <Text origin={opts.origin}>Thank you for your feedback!</Text>
        </StyledAlert>
      ) : (
        <form onSubmit={submitFeedback}>
          <StyledGrid
            container
            direction="row"
            alignItems="center"
            justify={opts.origin === "WIZARD" ? "center" : "flex-end"}
          >
            <MuiGrid item xs={opts.origin === "WIZARD" ? 6 : 12}>
              <Text origin={opts.origin}>{survey.prompt}</Text>
            </MuiGrid>
            <MuiGrid
              item
              xs={opts.origin === "WIZARD" ? 6 : 12}
              style={{ display: "flex", justifyContent: "space-around" }}
            >
              <EmojiRatings
                size={opts.origin === "ANYTIME" ? "large" : "medium"}
                ratings={survey.ratingLabels}
                setRating={setSelected}
              />
            </MuiGrid>
            {selected !== null && (
              <>
                {opts.origin === "ANYTIME" && (
                  <MuiGrid item xs={12} style={{ marginTop: "32px", marginBottom: "16px" }}>
                    <Select
                      name="anytimeSelect"
                      label={survey.feedbackTypeLabel}
                      options={survey.feedbackType}
                      onChange={setType}
                    />
                  </MuiGrid>
                )}
                <MuiGrid item xs={12}>
                  <TextField
                    multiline
                    placeholder={survey.freeformPrompt}
                    label={opts.origin === "ANYTIME" ? survey.freeformLabel : null}
                    value={feedback}
                    helperText={`${feedback.trim().length} / ${maxLength}`}
                    error={error}
                    onChange={e => {
                      setFeedback(e.target.value);
                      setError(e.target.value.trim().length > maxLength);
                    }}
                    {...textFieldProps}
                  />
                </MuiGrid>
                <MuiGrid item xs={4}>
                  <Button
                    type="submit"
                    text="Submit"
                    variant={opts.origin === "WIZARD" ? "secondary" : "primary"}
                    disabled={error}
                  />
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
