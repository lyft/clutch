---
title: Workflow Feedback Collection
{{ .EditURL }}
---

Clutch has support for collecting and storing feedback on your workflows through:
- An in-line feedback placement on workflows that use the [Wizard component](https://clutch.sh/docs/development/frontend/overview/#clutch-shwizard), allowing users to leave feedback directly in the last step of the workflow operation.
- A feedback placement on the app header, allowing users to leave feedback anytime. Users can provide general feedback or specify a feedback type in their submission, so this placement is also useful for workflows that don’t use the Wizard.

## Components
The Clutch feedback framework consists of three parts: frontend, backend server, and [feature flags](https://clutch.sh/docs/development/feature-flags). Below are the components responsible for enabling feedback collection in Clutch. **Registration of all the components is required**.

| Name | Description |
| --- | --- |
| `clutch.module.feedback` | Provides the `getSurveys` and `submitFeedback` endpoints. The former returns the survey questions for the feedback placements and the latter calls the feedback service. |
| `clutch.module.featureflag` | Handles the feature flag configuration to enable the feedback features. |
| `clutch.service.db.postgres` | Connects to the Postgres database.  |
| `clutch.service.feedback` | Stores the feedback submissions in the database.  |

### Configuration
Below is an example [clutch-config](https://clutch.sh/docs/configuration/) for enabling feedback collection.

```yaml title="backend/clutch-config.yaml"
...
modules:
  ...
  - name: clutch.module.featureflag
    typed_config:
      "@type": types.google.com/clutch.config.module.featureflag.v1.Config
      simple:
        flags:
          # enables the feedback placement on all workflows that use the Wizard
          npsWizard: true
          # enables the feedback placement on the app header
          npsHeader: true
  - name: clutch.module.feedback
    typed_config:
      "@type": types.google.com/clutch.config.module.feedback.v1.Config
      origins:
        - origin: WIZARD
          survey:
            prompt: Rate the workflow
            freeform_prompt: What would you recommend to improve the workflow?
            rating_labels:
              - emoji: SAD
                label: bad
              - emoji: NEUTRAL
                label: ok
              - emoji: HAPPY
                label: great
        - origin: HEADER
          survey:
            prompt: Rate your experience using Clutch
            freeform_prompt: What's going well? What could be better?
            rating_labels:
              - emoji: SAD
                label: bad
              - emoji: NEUTRAL
                label: ok
              - emoji: HAPPY
                label: great
services:
  ...
  - name: clutch.service.db.postgres
    ...
  - name: clutch.service.feedback
```

#### Customization
The survey questions for each feedback placement are driven by the clutch-config. As seen in the example above, provide the question that you want the user to rate as well as the label (which will be used as the tooltip text for corresponding emoji option) for each emoji rating.

## Backend

### Feedback Service
The feedback service has two responsibilties: normalize the feedback rating to a score out a 100 (described below) and save the feedback to the database.

#### Rating System
The Clutch feedback framework currently ships with a three-point emoji rating scale, which gets normalized to a raw score out of a 100 (sad emoji: 30, neutral emoji: 70, happy emoji: 100). There are two main benefits:
1. By normalizing the rating to a score out of 100, existing rating scales can be extended or new rating scales can be added without needing to change the scores of other rating scales. Importantly, this also means that prior feedback submission scores can still be combined with the data set of the new or updated rating scales.
1. By normalizing the rating to a score out of 100, the feedback can be easily computed to generate NPS (Net Promoter Score) or CSAT (Customer Satisfaction Score).

#### Feedback Submission
Below is an example of what gets sent in the `submitFeedback` API request.

```json
{
    "userId": "foo@example.com",
    "id": "z2z749z4-z067-443z-zz95-9796511z7zzf",
    "feedback": {
        "ratingLabel": "great",
        "ratingScale": {
            "emoji": 3
        },
        "feedbackType": "/k8s/pod/delete",
        "freeformResponse": "this is great!"
    },
    "metadata": {
        "origin": 2,
        "userSubmitted": true,
        "survey": {
            "prompt": "Rate the workflow",
            "freeformPrompt": "What would you recommend to improve the workflow?",
            "ratingLabels": [
                {
                    "emoji": "SAD",
                    "label": "bad"
                },
                {
                    "emoji": "NEUTRAL",
                    "label": "ok"
                },
                {
                    "emoji": "HAPPY",
                    "label": "great"
                }
            ]
        },
        "urlSearchParams": "?q=cluster%2Fnamespace%2Fpod"
    }
}
```

Some pieces to highlight:
- `feedbackType`: the area of feedback for the submission (i.e. url path)
- `ratingScale`: the scale presented to the user as well as the scale option the user selected. This is then normalized to a score out of 100.
- `metadata`: additional contextual information on the feedback submission, such as what was the survey question, what was the origin (i.e. Wizard or header placement), the url search params (to know what the user was looking up), and whether the feedback was formally submitted (described below).

In order to collect as much feedback data as possible, a `submitSuvey` API call is made every time an emoji rating is selected, even if the user did not click the submit button. Because of this, the feedback service uses the `client_id` from the feedback submission to continue updating the user’s feedback throughout the session to ensure we store the most recent submission. As part of this, there is a `userSubmitted` field in the feedback metadata to know whether the submission was formally submitted.
