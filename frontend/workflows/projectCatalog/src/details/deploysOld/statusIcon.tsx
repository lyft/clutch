import React from "react";
import CheckCircleOutlineIcon from "@material-ui/icons/CheckCircleOutline";
import ErrorOutlineIcon from "@material-ui/icons/ErrorOutline";
import HighlightOffIcon from "@material-ui/icons/HighlightOff";
import PlayCircleOutlineIcon from "@material-ui/icons/PlayCircleOutline";
import QueueIcon from "@material-ui/icons/Queue";
import SkipNextIcon from "@material-ui/icons/SkipNext";
import TimerIcon from "@material-ui/icons/Timer";
import WarningIcon from "@material-ui/icons/Warning";

import type { Statuses } from "./types";

const StatusIcon = (status: Statuses) => {
  switch (status) {
    case "WAITING":
      return <TimerIcon style={{ color: "#B09027" }} />;
    case "RUNNING":
      return <PlayCircleOutlineIcon style={{ color: "#3548D4" }} />;
    case "SUCCESS":
      return <CheckCircleOutlineIcon style={{ color: "#40A05A" }} />;
    case "FAILURE":
      return <ErrorOutlineIcon style={{ color: "#C2302E" }} />;
    case "ABORTED":
      return <HighlightOffIcon style={{ color: "#D87313" }} />;
    case "SKIPPED":
      return <SkipNextIcon style={{ color: "#9192A1" }} />;
    case "QUEUED":
      return <QueueIcon />;
    case "UNKNOWN":
    default:
      return <WarningIcon style={{ color: "#D87313" }} />;
  }
};

export default StatusIcon;
