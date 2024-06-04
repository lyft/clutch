import React from "react";
import type { DateTimePickerProps as MuiDateTimePickerProps } from "@mui/x-date-pickers";
import { DateTimePicker as MuiDateTimePicker } from "@mui/x-date-pickers";
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

export interface DateTimePickerProps
  extends Pick<
    MuiDateTimePickerProps<Date, Date>,
    "disabled" | "value" | "onChange" | "label" | "minDate" | "maxDate"
  > {
  allowEmpty?: boolean;
  error?: boolean;
  helperText?: string;
}

const DateTimePicker = ({
  onChange,
  allowEmpty = false,
  error = false,
  helperText = "",
  minDate = null,
  maxDate = null,
  ...props
}: DateTimePickerProps) => (
  <LocalizationProvider dateAdapter={AdapterDayjs}>
    <MuiDateTimePicker
      renderInput={inputProps => (
        <PaddedTextField {...inputProps} error={error} helperText={helperText} />
      )}
      onChange={(value: Dayjs | null) => {
        const isDateValid = value && value.isValid();
        if (allowEmpty || (!allowEmpty && isDateValid)) {
          const dateValue = isDateValid ? value.toDate() : null;
          onChange(dateValue);
        }
      }}
      {...props}
    />
  </LocalizationProvider>
);

export default DateTimePicker;
