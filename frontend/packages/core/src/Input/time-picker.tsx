import React from "react";
import type { TimePickerProps as MuiTimePickerProps } from "@mui/lab";
import { TimePicker as MuiTimePicker } from "@mui/x-date-pickers";
import { AdapterDayjs } from "@mui/x-date-pickers/AdapterDayjs";
import { LocalizationProvider } from "@mui/x-date-pickers/LocalizationProvider";
import type { Dayjs } from "dayjs";

import styled from "../styled";

import { TextField } from "./text-field";

const PaddedTextField = styled(TextField)({
  // This is required as TextField intentionally unsets the right padding for
  // end adornment styles since material introduced it in v5
  // Clutch has TextFields with end adornments that are end aligned (e.g. resolvers).
  ".MuiInputBase-adornedEnd": {
    paddingRight: "14px",
  },
});

export interface TimePickerProps
  extends Pick<
    MuiTimePickerProps,
    "disabled" | "value" | "onChange" | "label" | "PaperProps" | "PopperProps"
  > {}

const TimePicker = ({ onChange, ...props }: TimePickerProps) => (
  <LocalizationProvider dateAdapter={AdapterDayjs}>
    <MuiTimePicker
      renderInput={inputProps => <PaddedTextField {...inputProps} />}
      onChange={(value: Dayjs | null) => {
        if (value && value.isValid()) {
          onChange(value.toDate());
        }
      }}
      {...props}
    />
  </LocalizationProvider>
);

export default TimePicker;
