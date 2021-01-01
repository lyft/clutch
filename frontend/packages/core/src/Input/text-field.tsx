import * as React from "react";
import styled from "@emotion/styled";
import type {
  InputProps as MuiInputProps,
  StandardTextFieldProps as MuiStandardTextFieldProps,
} from "@material-ui/core";
import { IconButton as MuiIconButton, TextField as MuiTextField } from "@material-ui/core";
import ErrorIcon from "@material-ui/icons/Error";

const KEY_ENTER = 13;

const BaseTextField = ({ InputProps, InputLabelProps, ...props }: MuiStandardTextFieldProps) => (
  <MuiTextField
    InputLabelProps={{ ...InputLabelProps, shrink: true }}
    InputProps={{ ...InputProps, disableUnderline: true }}
    fullWidth
    {...props}
  />
);

const StyledTextField = styled(BaseTextField)({
  ".MuiInputLabel-root": {
    fontSize: "14px",
    fontWeight: 500,
    transform: "scale(1)",
    marginLeft: "2px",
  },

  ".MuiInputLabel-root, .MuiInputLabel-root.Mui-focused": {
    color: "rgba(13, 16, 48, 0.6)",
  },

  ".MuiInputLabel-root.Mui-disabled": {
    color: "rgba(13, 16, 48, 0.38)",
  },

  ".MuiInputLabel-root.Mui-error": {
    color: "#db3615",
  },

  ".MuiInputBase-root": {
    border: "1px solid rgba(13, 16, 48, 0.38)",
    borderRadius: "4px",
    fontSize: "16px",
    color: "#0D1030"
  },

  "label + .MuiInput-formControl": {
    marginTop: "20px",
  },

  ".MuiInputBase-root.Mui-focused": {
    borderColor: "#3548d4",
    ".MuiButtonBase-root": {
      backgroundColor: "#3548D4",
      "*": {
        color: "#FFFFFF",
      },
    },
  },

  ".MuiInputBase-root.Mui-disabled": {
    backgroundColor: "rgba(13, 16, 48, 0.12)",
  },

  ".MuiInput-input": {
    padding: "14px 16px",
    height: "20px",
  },

  ".MuiInput-input::placeholder": {
    color: "rgba(13, 16, 48, 0.38)",
    opacity: 1,
  },

  ".MuiFormHelperText-root": {
    alignItems: "center",
    display: "flex",
    position: "relative",
    fontSize: "12px",
    marginTop: "7px",
    lineHeight: "16px",
  },

  ".MuiFormHelperText-root.Mui-error": {
    color: "#db3615",
  },

  ".MuiInputBase-root.Mui-error": {
    borderColor: "#db3615",
  },

  ".MuiFormHelperText-root > svg": {
    height: "16px",
    width: "16px",
    marginRight: "4px",
  },
});

const IconButton = styled(MuiIconButton)({
  borderRadius: "0",
  backgroundColor: "#E7E7EA",
  borderBottomRightRadius: "3px",
  borderTopRightRadius: "3px",
  "&:hover": {
    backgroundColor: "#2D3DB4",
    "*": {
      color: "#FFFFFF",
    },
  },
  "&:active": {
    backgroundColor: "#2938A5",
    "*": {
      color: "#FFFFFF",
    },
  },
});

export interface TextFieldProps
  extends Pick<
      MuiStandardTextFieldProps,
      | "defaultValue"
      | "disabled"
      | "error"
      | "helperText"
      | "id"
      | "inputRef"
      | "label"
      | "multiline"
      | "name"
      | "onChange"
      | "onFocus"
      | "onKeyDown"
      | "placeholder"
      | "required"
      | "type"
      | "value"
    >,
    Pick<MuiInputProps, "readOnly" | "endAdornment"> {
  onReturn?: () => void;
}

export const TextField = ({
  onChange,
  onReturn,
  error,
  helperText,
  readOnly,
  endAdornment,
  ...props
}: TextFieldProps) => {
  const onKeyDown = (
    e: React.KeyboardEvent<HTMLDivElement | HTMLTextAreaElement | HTMLInputElement>
  ) => {
    if (onChange !== undefined) {
      onChange(e as React.ChangeEvent<any>);
    }
    if (e.keyCode === KEY_ENTER && onReturn) {
      onReturn();
    }
  };

  let helpText = helperText;

  // Prepend a '!' icon to helperText displayed below the form if the form is in an error state.
  if (error) {
    helpText = (
      <>
        <ErrorIcon />
        {helperText}
      </>
    );
  }

  return (
    <StyledTextField
      onKeyDown={e => onKeyDown(e)}
      onFocus={onChange}
      onBlur={onChange}
      error={error}
      helperText={helpText}
      InputProps={{
        readOnly,
        endAdornment: endAdornment && <IconButton type="submit">{endAdornment}</IconButton>,
      }}
      {...props}
    />
  );
};

export default TextField;
