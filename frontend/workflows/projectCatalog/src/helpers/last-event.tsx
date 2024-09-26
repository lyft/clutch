import React from "react";
import { Grid, TimeAgo as EventTime, Typography, TypographyProps } from "@clutch-sh/core";
import { faClock } from "@fortawesome/free-regular-svg-icons";
import { FontAwesomeIcon, FontAwesomeIconProps } from "@fortawesome/react-fontawesome";
import type { GridProps } from "@mui/material";

interface Size {
  variant: TypographyProps["variant"];
  icon: FontAwesomeIconProps["size"];
  align: GridProps["alignItems"];
}

interface SizeMap {
  [key: string]: Size;
}

const SIZE_MAP: SizeMap = {
  small: {
    variant: "body4",
    icon: "sm",
    align: "center",
  },
  medium: {
    variant: "body3",
    icon: "1x",
    align: "flex-end",
  },
  large: {
    variant: "body2",
    icon: "lg",
    align: "flex-end",
  },
};

interface LastEventProps {
  time: number;
  size?: keyof SizeMap;
}

const LastEvent = ({ time, size = "small", ...props }: LastEventProps) => {
  return time ? (
    <Grid item>
      <Grid container spacing={0.5} alignItems={SIZE_MAP[size].align}>
        <Grid item>
          <FontAwesomeIcon icon={faClock} size={SIZE_MAP[size].icon} />
        </Grid>
        <Grid item>
          <Typography variant={SIZE_MAP[size].variant}>
            <EventTime date={time} {...props} /> ago
          </Typography>
        </Grid>
      </Grid>
    </Grid>
  ) : null;
};

export default LastEvent;
