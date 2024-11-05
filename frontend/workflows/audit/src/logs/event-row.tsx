import React from "react";
import ReactJson from "react-json-view";
import type { clutch as IClutch } from "@clutch-sh/api";
import type { ClutchError } from "@clutch-sh/core";
import { client, IconButton, Popper, PopperItem, styled, TableRow } from "@clutch-sh/core";
import CheckIcon from "@mui/icons-material/Check";
import DownloadIcon from "@mui/icons-material/Download";
import FileCopyIcon from "@mui/icons-material/FileCopyOutlined";
import KeyboardArrowDownIcon from "@mui/icons-material/KeyboardArrowDown";
import KeyboardArrowUpIcon from "@mui/icons-material/KeyboardArrowUp";
import MoreVertIcon from "@mui/icons-material/MoreVert";
import OpenInNewIcon from "@mui/icons-material/OpenInNew";
import ShareIcon from "@mui/icons-material/Share";
import { Stack, Theme } from "@mui/material";
import FileSaver from "file-saver";

const ENDPOINT = "/v1/audit/getEvents";
const COLUMN_COUNT = 6;
const MonospaceText = styled("div")(({ theme }: { theme: Theme }) => ({
  fontFamily: "monospace",
  padding: theme.spacing("sm"),
  border: "1px solid lightgray",
  borderRadius: "8px",
  background: theme.palette.secondary[200],
  overflow: "hidden",
  textOverflow: "ellipsis",
  whiteSpace: "pre",
  maxWidth: "400px",
}));

const ReactJsonWrapper = styled("div")(({ theme }: { theme: Theme }) => ({
  padding: theme.spacing("sm"),
}));

interface EventRowAction {
  icon: React.ReactElement;
  name: string;
  onClick: () => void;
  disabled?: boolean;
}

interface EventRowProps {
  event: IClutch.audit.v1.IEvent;
  detailsPathPrefix?: string;
  downloadPrefix?: string;
}

const EventRow = ({ event, detailsPathPrefix, downloadPrefix }: EventRowProps) => {
  const [open, setOpen] = React.useState<boolean>(false);
  const anchorRef = React.useRef(null);
  const [showActions, setShowActions] = React.useState<boolean>(false);
  const [shareClicked, setShareClicked] = React.useState<boolean>(false);
  const date = new Date(String(event.occurredAt)).toLocaleString();
  const requestBody = { ...event.event.requestMetadata.body };
  delete requestBody["@type"];
  const { methodName, serviceName, username } = event.event;

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
  ] as EventRowAction[];

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
          icon: shareClicked ? <CheckIcon /> : <ShareIcon />,
          name: shareClicked ? "Copied" : "Share",
          onClick: () => {
            navigator.clipboard.writeText(url);
            setShareClicked(true);
            setTimeout(() => {
              setShareClicked(false);
            }, 1000);
          },
          disabled: shareClicked,
        },
      ];
    } catch {
      // eslint-disable-next-line no-console
      console.warn(
        `invalid event URL: ${window.location.protocol}//${window.location.host}/${detailsPathPrefix}/${event.id}`
      );
    }
  }

  return (
    <>
      <TableRow>
        <>{date}</>
        <>{methodName}</>
        <>{serviceName}</>
        <MonospaceText>{JSON.stringify(requestBody, null, 1)}</MonospaceText>
        <>{username}</>
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
        <TableRow colSpan={COLUMN_COUNT}>
          <ReactJsonWrapper>
            <ReactJson
              src={event}
              name={null}
              groupArraysAfterLength={5}
              displayDataTypes={false}
              collapsed={2}
            />
          </ReactJsonWrapper>
        </TableRow>
      )}
      <Popper open={showActions} anchorRef={anchorRef} onClickAway={() => setShowActions(false)}>
        {actions.map(a => (
          <PopperItem key={a.name} icon={a.icon} onClick={a.onClick} disabled={a.disabled || false}>
            {a.name}
          </PopperItem>
        ))}
      </Popper>
    </>
  );
};

interface EventRowsProps extends Pick<EventRowProps, "detailsPathPrefix" | "downloadPrefix"> {
  startTime: Date;
  endTime: Date;
  rangeKey: string;
  onFetch: () => void;
  onSuccess: () => void;
  onError: (e: ClutchError) => void;
}

const EventRows = ({
  detailsPathPrefix,
  downloadPrefix,
  startTime,
  endTime,
  rangeKey,
  onFetch,
  onSuccess,
  onError,
}: EventRowsProps) => {
  const [isLoading, setIsLoading] = React.useState<boolean>(false);
  const [hasError, setHasError] = React.useState<boolean>(false);
  const [pageToken, setPageToken] = React.useState<string>("0");
  const [events, setEvents] = React.useState<IClutch.audit.v1.IEvent[]>([]);

  const containerRef = React.useRef(null);

  const fetch = (page?: string) => {
    if (((page === undefined || page === "") && pageToken === "") || isLoading) {
      return;
    }
    const tkn = page || pageToken;
    onFetch();
    setIsLoading(true);
    const data = {
      range: { startTime, endTime },
      limit: 10,
      pageToken: tkn,
    } as IClutch.audit.v1.IGetEventsRequest;
    client
      .post(ENDPOINT, data)
      .then(resp => {
        const response = resp?.data as IClutch.audit.v1.GetEventsResponse;
        if (response?.events) {
          setEvents(evnts => [...evnts, ...response.events]);
          setPageToken(response.nextPageToken);
        }
        setHasError(false);
        onSuccess();
      })
      .catch((e: ClutchError) => {
        setHasError(true);
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
      if (entry.isIntersecting && !hasError) {
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
    // n.b. we explicitly pass in 0 here in addition to the setPageToken
    // call above since react won't pick up the state updates when
    // fetch is invoked.
    fetch("0");
  }, [rangeKey]);

  if (events.length <= 0) {
    return (
      <TableRow colSpan={COLUMN_COUNT}>
        <div style={{ textAlign: "center" }}>No Events Found</div>
      </TableRow>
    );
  }
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

export default EventRows;
