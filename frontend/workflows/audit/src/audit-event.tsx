import React from "react";
import ReactJson from "react-json-view";
import type { clutch as IClutch } from "@clutch-sh/api";
import { client, ClutchError, Error, Loadable, Typography, useParams } from "@clutch-sh/core";
import { Stack, useTheme } from "@mui/material";

import type { WorkflowProps } from ".";

const ENDPOINT = "/v1/audit/getEvent";

const AuditEvent: React.FC<WorkflowProps> = ({ heading }) => {
  const theme = useTheme();
  const params = useParams();
  const [isLoading, setIsLoading] = React.useState<boolean>(true);
  const [event, setEvent] = React.useState<IClutch.audit.v1.IEvent>();
  const [error, setError] = React.useState<ClutchError | null>(null);

  const fetch = () => {
    setIsLoading(true);
    const data = {
      eventId: parseInt(params.id, 10),
    } as IClutch.audit.v1.IGetEventRequest;
    client
      .post(ENDPOINT, data)
      .then(resp => {
        const eventResponse = resp?.data as IClutch.audit.v1.GetEventResponse;
        if (eventResponse?.event) {
          setEvent(eventResponse.event);
        }
        setError(null);
      })
      .catch((e: ClutchError) => {
        setError(e);
      })
      .finally(() => setIsLoading(false));
  };

  React.useEffect(() => fetch(), []);

  return (
    <Stack spacing={2} direction="column" style={{ padding: theme.clutch.layout.gutter }}>
      {!theme.clutch.useWorkflowLayout && <Typography variant="h2">{heading}</Typography>}
      <Loadable isLoading={isLoading}>
        {error && <Error subject={error} />}
        <ReactJson
          src={event}
          name={null}
          groupArraysAfterLength={10}
          displayDataTypes={false}
          collapsed={3}
        />
      </Loadable>
    </Stack>
  );
};

export default AuditEvent;
