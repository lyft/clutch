import React from "react";
import { Button, DateTimePicker, IconButton } from "@clutch-sh/core";
import AccessTimeIcon from "@mui/icons-material/AccessTime";
import { Menu, MenuItem, Stack } from "@mui/material";

const QUICK_TIME_OPTIONS = [
  { label: "Last 1m", value: 1 },
  { label: "Last 5m", value: 5 },
  { label: "Last 30m", value: 30 },
  { label: "Last 1hr", value: 60 },
  { label: "Last 2hr", value: 120 },
  { label: "Last 6hr", value: 360 },
];

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

export { DateTimeRangeSelector, QUICK_TIME_OPTIONS };
