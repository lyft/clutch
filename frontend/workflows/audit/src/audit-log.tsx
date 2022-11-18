import React from "react";
import ReactJson from "react-json-view";
import type { clutch as IClutch } from "@clutch-sh/api";
import {
  Button,
  client,
  ClutchError,
  DateTimePicker,
  Error,
  IconButton,
  Popper,
  PopperItem,
  styled,
  Table,
  TableRow,
  Typography,
  useSearchParams,
} from "@clutch-sh/core";
import AccessTimeIcon from "@mui/icons-material/AccessTime";
import DownloadIcon from "@mui/icons-material/Download";
import FileCopyIcon from "@mui/icons-material/FileCopyOutlined";
import KeyboardArrowDownIcon from "@mui/icons-material/KeyboardArrowDown";
import KeyboardArrowUpIcon from "@mui/icons-material/KeyboardArrowUp";
import MoreVertIcon from "@mui/icons-material/MoreVert";
import OpenInNewIcon from "@mui/icons-material/OpenInNew";
import SearchIcon from "@mui/icons-material/Search";
import ShareIcon from "@mui/icons-material/Share";
import { Menu, MenuItem, Stack, useMediaQuery, useTheme } from "@mui/material";
import FileSaver from "file-saver";

import type { AuditLogProps } from ".";

const ENDPOINT = "/v1/audit/getEvents";
const QUICK_TIME_OPTIONS = [
  { label: "Last 1m", value: 1 },
  { label: "Last 5m", value: 5 },
  { label: "Last 30m", value: 30 },
  { label: "Last 1hr", value: 60 },
  { label: "Last 2hr", value: 120 },
  { label: "Last 6hr", value: 360 },
];

const MonospaceText = styled("div")({
  fontFamily: "monospace",
  padding: "8px",
  border: "1px solid lightgray",
  borderRadius: "8px",
  background: "#ddd9d9",
  overflow: "hidden",
  textOverflow: "ellipsis",
  whiteSpace: "pre",
  maxWidth: "400px",
});

interface EventRowProps {
  event: IClutch.audit.v1.IEvent;
  detailsPathPrefix?: string;
  downloadPrefix?: string;
}

const EventRow = ({ event, detailsPathPrefix, downloadPrefix }: EventRowProps) => {
  const [open, setOpen] = React.useState<boolean>(false);
  const anchorRef = React.useRef(null);
  const [showActions, setShowActions] = React.useState<boolean>(false);
  const date = new Date(String(event.occurredAt)).toLocaleString();
  const requestBody = { ...event.event.requestMetadata.body };
  delete requestBody["@type"];
  const method = event.event.methodName;
  const service = event.event.serviceName;

  let actions = [
    {
      icon: <FileCopyIcon />,
      name: "Copy",
      onClick: () => {
        const output = JSON.stringify(event, null, "\t");
        navigator.clipboard.writeText(output);
      },
    },
    {
      icon: <DownloadIcon />,
      name: "Download",
      onClick: () => {
        const output = new Blob([JSON.stringify(event, null, "\t")]);
        const prefix = downloadPrefix || "clutch_audit_event";
        FileSaver.saveAs(output, `${prefix}_${event.id}_${Date.now()}.json`);
      },
    },
  ];
  if (detailsPathPrefix) {
    // if a URL cannot be created for the event details page gracefully fallback to showing nothing
    try {
      const url = new URL(
        `${detailsPathPrefix}/${event.id}`,
        `${window.location.protocol}//${window.location.host}`
      ).toString();
      actions = [
        {
          icon: <OpenInNewIcon />,
          name: "Open",
          onClick: () => window.open(url, "_blank"),
        },
        ...actions,
        {
          icon: <ShareIcon />,
          name: "Share",
          onClick: () => navigator.clipboard.writeText(url),
        },
      ];
    } catch {}
  }

  return (
    <>
      <TableRow>
        <>{date}</>
        <>{method}</>
        <>{service}</>
        <MonospaceText>{JSON.stringify(requestBody, null, 1)}</MonospaceText>
        <>{event.event.username}</>
        <Stack direction="row">
          <IconButton variant="neutral" onClick={() => setOpen(o => !o)}>
            {open ? <KeyboardArrowUpIcon /> : <KeyboardArrowDownIcon />}
          </IconButton>
          <IconButton variant="neutral" onClick={() => setShowActions(o => !o)} ref={anchorRef}>
            <MoreVertIcon />
          </IconButton>
        </Stack>
      </TableRow>
      {open && (
        <TableRow colSpan={6}>
          <div style={{ padding: "8px" }}>
            <ReactJson
              src={event}
              name={null}
              groupArraysAfterLength={5}
              displayDataTypes={false}
              collapsed={2}
            />
          </div>
        </TableRow>
      )}
      <Popper open={showActions} anchorRef={anchorRef} onClickAway={() => setShowActions(false)}>
        {actions.map(a => (
          <PopperItem key={a.name} icon={a.icon} onClick={a.onClick}>
            {a.name}
          </PopperItem>
        ))}
      </Popper>
    </>
  );
};

interface DateTimeRangeSelectorProps {
  shrink: boolean;
  disabled: boolean;
  start?: Date;
  end?: Date;
  onStartChange: (start: Date) => void;
  onEndChange: (end: Date) => void;
  onQuickSelect: (start: Date, end: Date) => void;
}

const DateTimeRangeSelector = ({
  shrink,
  disabled,
  start,
  end,
  onStartChange,
  onEndChange,
  onQuickSelect,
}: DateTimeRangeSelectorProps) => {
  const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null);
  const open = Boolean(anchorEl);

  const calculateQuickTime = (value: number) => {
    const currentTime = new Date();
    setAnchorEl(null);
    onQuickSelect(new Date(currentTime.getTime() - value * 60 * 1000), currentTime);
  };
  return (
    <>
      <Stack direction={shrink ? "column" : "row"} spacing={1}>
        <DateTimePicker
          label="Start Time"
          value={start}
          onChange={(newValue: Date) => {
            onStartChange(newValue);
          }}
          disabled={disabled}
        />
        <DateTimePicker
          label="End Time"
          value={end}
          onChange={(newValue: Date) => {
            onEndChange(newValue);
          }}
          disabled={disabled}
        />
        {shrink ? (
          <Button
            text="Quick Time Select"
            onClick={e => setAnchorEl(e.currentTarget)}
            variant="neutral"
          />
        ) : (
          <IconButton
            onClick={e => setAnchorEl(e.currentTarget)}
            variant="neutral"
            disabled={disabled}
          >
            <AccessTimeIcon />
          </IconButton>
        )}
      </Stack>
      <Menu anchorEl={anchorEl} open={open} onClose={() => setAnchorEl(null)}>
        {QUICK_TIME_OPTIONS.map(o => (
          <MenuItem key={o.label} onClick={() => calculateQuickTime(o.value)}>
            {o.label}
          </MenuItem>
        ))}
      </Menu>
    </>
  );
};

interface EventRowsProps extends Pick<EventRowProps, "detailsPathPrefix" | "downloadPrefix"> {
  startTime: Date;
  endTime: Date;
  key: string;
  onFetch: () => void;
  onSuccess: () => void;
  onError: (e: ClutchError) => void;
}

const EventRows = ({
  detailsPathPrefix,
  downloadPrefix,
  startTime,
  endTime,
  key,
  onFetch,
  onSuccess,
  onError,
}: EventRowsProps) => {
  const [isLoading, setIsLoading] = React.useState<boolean>(false);
  const [pageToken, setPageToken] = React.useState<string>("0");
  const [events, setEvents] = React.useState<IClutch.audit.v1.IEvent[]>([]);

  const containerRef = React.useRef(null);

  const fetch = () => {
    if (pageToken === "" || isLoading) {
      return;
    }
    onFetch();
    setIsLoading(true);
    const data = {
      range: { startTime, endTime },
      pageToken,
    } as IClutch.audit.v1.IGetEventsRequest;
    client
      .post(ENDPOINT, data)
      .then(resp => {
        const response = resp?.data as IClutch.audit.v1.GetEventsResponse;
        if (response?.events) {
          setEvents(evnts => [...evnts, ...response.events]);
          setPageToken(response.nextPageToken);
        }
        onSuccess();
      })
      .catch((e: ClutchError) => {
        onError(e);
      })
      .finally(() => setIsLoading(false));
  };

  const options = {
    root: null,
    rootMargin: "0px",
    threshold: 1.0,
  };
  React.useEffect(() => {
    const observer = new IntersectionObserver(entries => {
      const [entry] = entries;
      if (entry.isIntersecting) {
        fetch();
      }
    }, options);

    if (containerRef.current) {
      observer.observe(containerRef.current);
    }
    return () => {
      if (containerRef.current) {
        observer.unobserve(containerRef.current);
      }
    };
  }, [containerRef, options]);
  React.useEffect(() => {
    setPageToken("0");
    setEvents([]);
    fetch();
  }, [key]);

  return (
    <>
      {events.map((r, idx) => (
        <React.Fragment key={String(r.id)}>
          <EventRow
            event={r}
            detailsPathPrefix={detailsPathPrefix}
            downloadPrefix={downloadPrefix}
          />
          {idx === events.length - 2 && <div key="observer" ref={containerRef} />}
        </React.Fragment>
      ))}
    </>
  );
};

const AuditLog: React.FC<AuditLogProps> = ({ heading, detailsPathPrefix, downloadPrefix }) => {
  const [searchParams, setSearchParams] = useSearchParams();
  const [now] = React.useState<Date>(new Date());
  const [key, setKey] = React.useState<string>("");
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
  }, [key]);

  const theme = useTheme();
  const shrink = useMediaQuery(theme.breakpoints.down("md"));

  return (
    <Stack spacing={2} direction="column" style={{ padding: "32px", height: "100%" }}>
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
              setKey(`${startTime}-${endTime}`);
            }}
          />
          {shrink ? (
            <Button text="Search" onClick={() => setKey(`${startTime}-${endTime}`)} />
          ) : (
            <IconButton onClick={() => setKey(`${startTime}-${endTime}`)}>
              <SearchIcon />
            </IconButton>
          )}
        </Stack>
        {error && <Error subject={error} />}
      </Stack>
      <div style={{ display: "flex", justifyContent: "center", overflow: "auto" }}>
        {/* <Loadable isLoading={isLoading} variant="overlay"> */}
        <Table
          stickyHeader
          columns={["Timestamp", "Action", "Service", "Request Body", "User"]}
          actionsColumn
          overflow="break-word"
        >
          <EventRows
            detailsPathPrefix={detailsPathPrefix}
            downloadPrefix={downloadPrefix}
            startTime={startTime}
            endTime={endTime}
            key={key}
            onFetch={() => setIsLoading(true)}
            onSuccess={() => setIsLoading(false)}
            onError={e => {
              setIsLoading(false);
              setError(e);
            }}
          />
        </Table>
        {/* </Loadable> */}
      </div>
    </Stack>
  );
};

export default AuditLog;
