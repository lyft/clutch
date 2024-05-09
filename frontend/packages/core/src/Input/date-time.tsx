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
  extends Pick<MuiDateTimePickerProps<Date, Date>, "disabled" | "value" | "onChange" | "label"> {
  allowEmpty: boolean;
}

const DateTimePicker = ({ onChange, allowEmpty = false, ...props }: DateTimePickerProps) => (
  <LocalizationProvider dateAdapter={AdapterDayjs}>
    <MuiDateTimePicker
      renderInput={inputProps => <PaddedTextField {...inputProps} />}
      onChange={(value: Dayjs | null) => {
        if (!allowEmpty && value && value.isValid()) {
          const dateValue = value ? value.toDate() : null;
          onChange(dateValue);
        }

        if (allowEmpty && !value) {
          onChange(null);
        }
      }}
      {...props}
    />
  </LocalizationProvider>
);

export default DateTimePicker;
