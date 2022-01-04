import React, { useState } from "react";
import { clutch as IClutch } from "@clutch-sh/api";
import styled from "@emotion/styled";
import { capitalize, Grid as MuiGrid } from "@material-ui/core";
import MuiSuccessIcon from "@material-ui/icons/CheckCircle";
import { debounce } from "lodash";
import { v4 as uuid } from "uuid";

import { userId } from "../AppLayout/user";
import { Button } from "../button";
import { Alert } from "../Feedback";
import { Select, TextField } from "../Input";
import { client } from "../Network";
import type { ClutchError } from "../Network/errors";
import { Typography } from "../typography";

import EmojiRatings, { Rating } from "./emojiRatings";

/** Interfaces */

type Origins = "WIZARD" | "ANYTIME";

interface FeedbackOptions {
  origin: Origins;
  onSubmit?: (submit: boolean) => void;
}

// TODO: (jslaughter) update with milestone 2 anytime typing
interface Survey extends IClutch.feedback.v1.ISurvey {
  freeformLabel?: string;
  feedbackTypeLabel?: string;
  feedbackType?: { label: string }[];
}

// Defaults in case of API failure
export const defaults: Survey = {
  prompt: "Rate Your Experience",
  freeformPrompt: "What would you recommend to improve this?",
  freeformLabel: "Do you have any thoughts you'd like to share?",
  feedbackTypeLabel: "Choose a type of feedback you want to submit",
  feedbackType: [{ label: "general" }, { label: "dash" }, { label: "workflows" }],
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

const StyledButton = styled(Button)<{ origin: Origins }>({}, props =>
  props.origin === "WIZARD"
    ? {
        fontSize: "14px",
        padding: "0 8px",
        height: "18px",
      }
    : {
        marginTop: "16px",
      }
);

const StyledTextField = styled(TextField)<{ origin: Origins }>(
  {
    marginTop: "15px",
  },
  props => ({
    ".MuiInputBase-root": {
      fontSize: props.origin === "WIZARD" ? "14px" : "16px",
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
const NPSFeedback = ({ origin = "ANYTIME", onSubmit }: FeedbackOptions) => {
  const [hasSubmit, setHasSubmit] = useState<boolean>(false);
  const [selected, setSelected] = useState<Rating>(null);
  const [freeformFeedback, setFreeformFeedback] = useState<string>("");
  const [error, setError] = useState<boolean>(false);
  const [survey, setSurvey] = useState<Survey>({});
  // const [survey, setSurvey] = useState<IClutch.feedback.v1.ISurvey>({});
  const [type, setType] = useState<string>(null);
  const [requestId, setRequestId] = useState<string>("");
  const maxLength = 180;
  const debounceTimer = 500;

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
    // let data: IClutch.feedback.v1.ISurvey = defaults;
    let data: Survey = defaults;

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
          feedbackType: type,
        },
        metadata: {
          origin: IClutch.feedback.v1.Origin[origin],
          userSubmitted: hasSubmit,
          survey,
          urlSearchParams: window.location.search,
        },
      });
    }
  }, [selected, hasSubmit]);

  // Form onSubmit handler
  const submitFeedback = e => {
    if (e) {
      e.preventDefault();
    }
    setHasSubmit(true);
    if (onSubmit) {
      onSubmit(true);
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
        justify={origin === "WIZARD" ? "center" : "flex-end"}
        style={{ padding: "16px" }}
      >
        <MuiGrid item xs={origin === "WIZARD" ? 6 : 12}>
          <Typography variant={origin === "WIZARD" ? "subtitle3" : "subtitle2"}>
            {survey.prompt}
          </Typography>
        </MuiGrid>
        <MuiGrid
          item
          xs={origin === "WIZARD" ? 6 : 12}
          style={{ display: "flex", justifyContent: "space-around" }}
        >
          <EmojiRatings
            ratings={survey.ratingLabels}
            setRating={setSelected}
            placement={origin === "WIZARD" ? "top" : "bottom"}
            size={origin === "WIZARD" ? "small" : "large"}
          />
        </MuiGrid>
        {selected !== null && (
          <>
            {origin === "ANYTIME" && (
              <MuiGrid item xs={12} style={{ marginTop: "32px", marginBottom: "16px" }}>
                <Select
                  name="anytimeSelect"
                  label={survey.feedbackTypeLabel}
                  options={survey.feedbackType.map(m => ({ label: capitalize(m.label) }))}
                  onChange={setType}
                />
              </MuiGrid>
            )}
            <MuiGrid item xs={12}>
              <StyledTextField
                multiline
                fullWidth
                origin={origin}
                placeholder={survey.freeformPrompt}
                label={origin === "ANYTIME" ? survey.freeformLabel : null}
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
                justifyContent: origin === "WIZARD" ? "center" : "flex-end",
              }}
            >
              <StyledButton
                origin={origin}
                type="submit"
                text="Submit"
                variant={origin === "WIZARD" ? "secondary" : "primary"}
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
