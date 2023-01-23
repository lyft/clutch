import * as React from "react";
import type { FieldValues, UseFormRegister } from "react-hook-form";
import styled from "@emotion/styled";
import ErrorIcon from "@mui/icons-material/Error";
import WarningIcon from "@mui/icons-material/Warning";
import type {
  InputProps as MuiInputProps,
  StandardTextFieldProps as MuiStandardTextFieldProps,
} from "@mui/material";
import {
  Autocomplete,
  Grid,
  IconButton as MuiIconButton,
  Popper as MuiPopper,
  TextField as MuiTextField,
  Typography,
} from "@mui/material";
import _ from "lodash";

const KEY_ENTER = 13;

const BaseTextField = React.forwardRef<HTMLInputElement, MuiStandardTextFieldProps>(
  ({ InputProps, InputLabelProps, ...props }: MuiStandardTextFieldProps, ref) => (
    <MuiTextField
      InputLabelProps={{ ...InputLabelProps, shrink: true }}
      InputProps={{ ...InputProps }}
      fullWidth
      {...props}
      {...{ ref }}
    />
  )
);

const StyledAutocomplete = styled(Autocomplete)({
  ".MuiOutlinedInput-root.MuiInputBase-sizeSmall": {
    padding: "unset",
  },
  height: "unset",
  ".MuiTextField-root > .MuiInputBase-root > .MuiInputBase-input": {
    height: "20px",

    "&.MuiAutocomplete-input": {
      padding: "14px 16px",
    },
  },
});

const StyledTextField = styled(BaseTextField)({
  "--error-color": "#DB3615",
  "--warning-color": "#FFCC80",
  height: "unset",

  ".MuiInputLabel-root": {
    "--label-color": "rgba(13, 16, 48, 0.6)",

    color: "var(--label-color)",

    "&.Mui-focused": {
      color: "var(--label-color)",
    },
    "&.Mui-error": {
      color: "var(--error-color)",
    },
  },

  "&.Mui-warning .MuiInputLabel-root:not(.Mui-error), .MuiFormHelperText-root": {
    color: "var(--warning-color)",
  },

  ".MuiInputBase-root": {
    "--input-border-width": "1px",
    "--input-default-color": "rgba(13, 16, 48, 0.38)",
    borderRadius: "4px",
    fontSize: "16px",
    backgroundColor: "#FFFFFF",

    "& fieldset": {
      borderColor: "var(--input-default-color)",
    },
    "&.Mui-focused fieldset": {
      borderColor: "#3548d4",
      borderWidth: "var(--input-border-width)",
    },
    "&.Mui-error fieldset": {
      borderColor: "var(--error-color)",
      borderWidth: "var(--input-border-width)",
    },
    "&.Mui-disabled fieldset": {
      backgroundColor: "rgba(13, 16, 48, 0.12)",
    },
    "& .MuiInputBase-input": {
      textOverflow: "ellipsis",
    },
    "> .MuiInputBase-input": {
      "--input-padding": "14px 16px",
      padding: "var(--input-padding)",
      height: "20px",

      "&.MuiAutocomplete-input": {
        padding: "var(--input-padding)",
      },

      "::placeholder": {
        color: "var(--input-default-color)",
        opacity: 1,
      },
    },
  },

  "&.Mui-warning .MuiInputBase-root:not(.Mui-error) fieldset": {
    "--input-border-width": "1px",
    borderColor: "var(--warning-color)",
    borderWidth: "var(--input-border-width)",
  },

  ".MuiInputBase-adornedEnd": {
    paddingRight: "unset",
  },

  ".MuiFormHelperText-root": {
    alignItems: "center",
    display: "flex",
    position: "relative",
    fontSize: "12px",
    marginTop: "7px",
    lineHeight: "16px",
    marginLeft: "0px",

    "&.Mui-error": {
      color: "var(--error-color)",
    },

    "> svg": {
      height: "16px",
      width: "16px",
      marginRight: "4px",
    },
  },
});

// popper containing the search result options
const Popper = styled(MuiPopper)({
  ".MuiPaper-root": {
    boxShadow: "0px 5px 15px rgba(53, 72, 212, 0.2)",

    "> .MuiAutocomplete-listbox": {
      "> .MuiAutocomplete-option": {
        height: "48px",
        padding: "0px",

        "&.Mui-focused": {
          background: "#ebedfb",
        },
      },
    },
  },
  ".MuiAutocomplete-noOptions": {
    fontSize: "14px",
    color: "#0d1030",
  },
});

// search's result options container
const ResultGrid = styled(Grid)({
  height: "inherit",
  padding: "12px 16px 12px 16px",
});

// search's result options
const ResultLabel = styled(Typography)({
  color: "#0d1030",
  fontSize: "14px",
});

const IconButton = styled(MuiIconButton)({
  borderRadius: "0",
  backgroundColor: "#3548D4",
  color: "#FFFFFF",
  borderBottomRightRadius: "3px",
  borderTopRightRadius: "3px",
  "&:hover": {
    backgroundColor: "#2D3DB4",
  },
  "&:active": {
    backgroundColor: "#2938A5",
  },
});

interface AutocompleteResultProps {
  id?: string;
  label: string;
}

const AutocompleteResult: React.FC<AutocompleteResultProps> = ({ id, label }) => (
  <ResultGrid container alignItems="center">
    <Grid item xs>
      <ResultLabel>{label || id}</ResultLabel>
    </Grid>
  </ResultGrid>
);

interface CustomInputProps {
  warning?: boolean;
}

export interface TextFieldProps
  extends Pick<
      MuiStandardTextFieldProps,
      | "defaultValue"
      | "disabled"
      | "error"
      | "fullWidth"
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
    Pick<CustomInputProps, "warning">,
    Pick<MuiInputProps, "readOnly" | "endAdornment"> {
  onReturn?: () => void;
  autocompleteCallback?: (v: string) => Promise<{ results: { id?: string; label: string }[] }>;
  formRegistration?: UseFormRegister<FieldValues>;
}

const TextFieldRef = (
  {
    onChange,
    onReturn,
    error,
    warning,
    helperText,
    readOnly,
    endAdornment,
    autocompleteCallback,
    defaultValue,
    value,
    fullWidth = true,
    name,
    required,
    formRegistration,
    inputRef,
    ...props
  }: TextFieldProps,
  ref
) => {
  const formValidation =
    formRegistration !== undefined ? formRegistration(name, { required }) : undefined;
  const changeCallback = onChange !== undefined ? onChange : e => {};
  const onKeyDown = (
    e: React.KeyboardEvent<HTMLDivElement | HTMLTextAreaElement | HTMLInputElement>
  ) => {
    if (formValidation !== undefined) {
      formValidation.onChange(e);
    }
    changeCallback(e as React.ChangeEvent<any>);
    if (e.keyCode === KEY_ENTER && onReturn && !error) {
      onReturn();
    }
  };

  let helpText = helperText;

  // Prepend a '!' icon to helperText displayed below the form if the form is in an error state.
  if ((error || warning) && helpText) {
    helpText = (
      <>
        {error && <ErrorIcon />}
        {!error && warning && <WarningIcon />}
        {helperText}
      </>
    );
  }

  // We maintain a defaultVal to prevent the value from changing from underneath
  // the component. This is required because autocomplete is uncontrolled.
  const [defaultVal] = React.useState<string>((defaultValue as string) || "");
  const [autoCompleteOptions, setAutoCompleteOptions] = React.useState<AutocompleteResultProps[]>(
    []
  );

  const isEmpty = (defaultValue === undefined || defaultValue === "") && value === "";
  const textFieldProps = {
    name,
    onFocus: e => {
      if (formValidation !== undefined) {
        formValidation.onChange(e);
      }
      changeCallback(e);
    },
    onBlur: e => {
      if (formValidation !== undefined) {
        formValidation.onBlur(e);
      }
      changeCallback(e);
    },
    error,
    helperText: helpText,
    InputProps: {
      onChange: e => {
        if (formValidation !== undefined) {
          formValidation.onChange(e);
        }
        changeCallback(e);
      },
      onKeyDown,
      readOnly,
      endAdornment: endAdornment ? (
        <IconButton type="submit" disabled={isEmpty} size="large">
          {endAdornment}
        </IconButton>
      ) : null,
    },
    inputRef: formValidation !== undefined ? formValidation.ref : inputRef,
  };

  const autoCompleteDebounce = React.useRef(
    _.debounce(val => {
      if (autocompleteCallback !== undefined) {
        autocompleteCallback(val)
          .then(data => {
            setAutoCompleteOptions(data.results);
          })
          .catch(err => {
            helpText = err;
          });
      }
    }, 250)
  ).current;
  if (autocompleteCallback !== undefined) {
    // TODO (mcutalo): support option.label in the renderOption
    return (
      <StyledAutocomplete
        freeSolo
        size="small"
        fullWidth={fullWidth}
        options={autoCompleteOptions}
        PopperComponent={Popper}
        getOptionLabel={(option: string | any) =>
          typeof option === "string" ? option : option?.id || option.label
        }
        onInputChange={(__, v) => autoCompleteDebounce(v)}
        renderOption={(otherProps, option: AutocompleteResultProps) => (
          <li className="MuiAutocomplete-option" {...otherProps}>
            <AutocompleteResult key={option.id} id={option.id} label={option.label} />
          </li>
        )}
        onSelectCapture={e => {
          if (formValidation !== undefined) {
            formValidation.onChange(e);
          }
          changeCallback(e as React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>);
        }}
        defaultValue={{ id: defaultVal, label: defaultVal }}
        value={value}
        renderInput={inputProps => (
          <StyledTextField
            {...inputProps}
            {...textFieldProps}
            InputProps={{
              ...textFieldProps.InputProps,
              ref: inputProps.InputProps.ref,
            }}
            {...props}
            {...{ ref }}
          />
        )}
        // This func is here for autocomplete. When the user clicks a choice in the dropdown
        // or presses enter, onChange is called, which will allow the user to submit their choice.
        // Without this, the user has to click their choice or press enter, then submit again once
        // the choice has updated. Note that this does not work if the `value` prop is being set
        // manually (as is the case in proj selector and proj catalog)
        // TODO: Make it work for all cases, not just the resolver and k8s dash.
        onChange={(_e, v) => {
          if (v && onReturn) {
            onReturn();
          }
        }}
      />
    );
  }

  return (
    <StyledTextField
      size="small"
      {...textFieldProps}
      defaultValue={defaultValue}
      value={value}
      onChange={onChange}
      className={warning ? "Mui-warning" : ""}
      {...props}
      {...{ ref }}
    />
  );
};

const TextField = React.forwardRef<HTMLDivElement, TextFieldProps>(TextFieldRef);

export default TextField;
