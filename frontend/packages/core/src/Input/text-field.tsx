import React from "react";
import type { InputLabelProps, TextFieldProps as MuiTextFieldProps } from "@material-ui/core";
import { TextField as MuiTextField } from "@material-ui/core";
import styled from "styled-components";

const KEY_ENTER = 13;
const StyledTextField = styled(MuiTextField)`
  ${({ theme, ...props }) => `
  display: flex;
  width: 100%;
  max-width: ${props["data-max-width"] || "500px"};
  margin: 15px;
  .MuiInputLabel-root {
    color: ${theme.palette.text.primary};
  }
  .MuiInput-underline:after {
    border-bottom: 2px solid ${theme.palette.accent.main};
  }
  `}
`;

export interface TextFieldProps {
  maxWidth?: string;
  onReturn?: () => void;
}

const TextField: React.FC<TextFieldProps & MuiTextFieldProps> = ({
  onChange,
  onReturn,
  maxWidth,
  placeholder,
  type,
  ...props
}) => {
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

  const inputLabelProps = {
    color: "secondary",
  } as Partial<InputLabelProps>;

  const placeholderInputTypes = ["date", "datetime-local", "month", "time", "week"];
  const hasInputPlaceholder = placeholderInputTypes.indexOf(type) !== -1;
  if (placeholder || hasInputPlaceholder) {
    inputLabelProps.shrink = true;
  }

  return (
    <StyledTextField
      data-max-width={maxWidth}
      color="secondary"
      InputLabelProps={inputLabelProps}
      onKeyDown={e => onKeyDown(e)}
      onFocus={onChange}
      onBlur={onChange}
      placeholder={placeholder}
      type={type}
      {...props}
    />
  );
};

export default TextField;
