import React from "react";
import {
  Button,
  ClutchError,
  Error,
  IconButton,
  styled,
  Table,
  Typography,
  useSearchParams,
} from "@clutch-sh/core";
import SearchIcon from "@mui/icons-material/Search";
import { CircularProgress, Stack, useMediaQuery, useTheme } from "@mui/material";

import type { AuditLogProps } from "..";

import { DateTimeRangeSelector, QUICK_TIME_OPTIONS } from "./date-selector";
import EventRows from "./event-row";

const RootContainer = styled(Stack)({
  padding: "32px",
  height: "100%",
});

const TableContainer = styled("div")({
  display: "flex",
  justifyContent: "center",
  overflow: "auto",
});

const LoadingContainer = styled("div")({
  height: "40px",
  width: "40px",
});

const LoadingSpinner = styled(CircularProgress)`
  color: #3548d4;
  position: absolute;
`;

const AuditLog: React.FC<AuditLogProps> = ({ heading, detailsPathPrefix, downloadPrefix }) => {
  const [searchParams, setSearchParams] = useSearchParams();
  const [now] = React.useState<Date>(new Date());
  const [timeRangeKey, setTimeRangeKey] = React.useState<string>("");
  const [isLoading, setIsLoading] = React.useState<boolean>(false);
  const [startTime, setStartTime] = React.useState<Date>(
    searchParams.get("start")
      ? new Date(searchParams.get("start"))
      : new Date(now.getTime() - QUICK_TIME_OPTIONS[0].value * 60 * 1000)
  );
  const [endTime, setEndTime] = React.useState<Date>(
    searchParams.get("end") ? new Date(searchParams.get("end")) : now
  );
  const [error, setError] = React.useState<ClutchError | null>(null);

  // n.b. on a time range change, attempt to update the search params and fail silently
  // as this is a nice to have.
  React.useEffect(() => {
    try {
      setSearchParams(
        {
          start: startTime.toISOString(),
          end: endTime.toISOString(),
        },
        { replace: true }
      );
    } catch {}
  }, [timeRangeKey]);

  const theme = useTheme();
  const shrink = useMediaQuery(theme.breakpoints.down("md"));

  return (
    <RootContainer spacing={2} direction="column">
      <Typography variant="h2">{heading}</Typography>
      <Stack direction="column" spacing={2}>
        <Stack
          direction={shrink ? "column" : "row"}
          spacing={1}
          sx={{
            alignSelf: shrink ? "center" : "flex-end",
            width: shrink ? "100%" : "inherit",
          }}
        >
          {isLoading && (
            <LoadingContainer>
              <LoadingSpinner />
            </LoadingContainer>
          )}
          <DateTimeRangeSelector
            shrink={shrink}
            disabled={isLoading}
            start={startTime}
            end={endTime}
            onStartChange={setStartTime}
            onEndChange={setEndTime}
            onQuickSelect={(start, end) => {
              setStartTime(start);
              setEndTime(end);
              setTimeRangeKey(`${startTime}-${endTime}`);
            }}
          />
          {shrink ? (
            <Button text="Search" onClick={() => setTimeRangeKey(`${startTime}-${endTime}`)} />
          ) : (
            <IconButton onClick={() => setTimeRangeKey(`${startTime}-${endTime}`)}>
              <SearchIcon />
            </IconButton>
          )}
        </Stack>
        {error && <Error subject={error} />}
      </Stack>
      <TableContainer>
        <Table
          stickyHeader
          columns={["Timestamp", "Action", "Service", "Request Body", "User"]}
          actionsColumn
          overflow={shrink ? "scroll" : "break-word"}
        >
          <EventRows
            detailsPathPrefix={detailsPathPrefix}
            downloadPrefix={downloadPrefix}
            startTime={startTime}
            endTime={endTime}
            rangeKey={timeRangeKey}
            onFetch={() => setIsLoading(true)}
            onSuccess={() => {
              setIsLoading(false);
              setError(null);
            }}
            onError={e => {
              setIsLoading(false);
              setError(e);
            }}
          />
        </Table>
      </TableContainer>
    </RootContainer>
  );
};

export default AuditLog;
