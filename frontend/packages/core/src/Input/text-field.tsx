import React from "react";
import type { TextFieldProps as MuiTextFieldProps } from "@material-ui/core";
import { TextField as MuiTextField } from "@material-ui/core";
import styled from "styled-components";

const KEY_ENTER = 13;
const StyledTextField = styled(MuiTextField)`
  ${({ theme, ...props }) => `
  display: flex;
  width: 100%;
  max-width: ${props["data-max-width"] || "500px"};
  margin: 5px;
  .MuiInputLabel-root {
    color: ${theme.palette.text.primary};
  }
  .MuiInput-underline:after {
    border-bottom: 2px solid ${theme.palette.accent.main};
  }
  `}
`;

interface TextFieldProps {
  maxWidth?: string;
  onReturn?: () => void;
}

const TextField: React.FC<TextFieldProps & MuiTextFieldProps> = ({
  onChange,
  onReturn,
  maxWidth,
  ...props
}) => {
  const onKeyDown = (
    e: React.KeyboardEvent<HTMLDivElement | HTMLTextAreaElement | HTMLInputElement>
  ) => {
    onChange(e as React.ChangeEvent<any>);
    if (e.keyCode === KEY_ENTER && onReturn) {
      onReturn();
    }
  };

  return (
    <StyledTextField
      data-max-width={maxWidth}
      color="secondary"
      InputLabelProps={{ color: "secondary" }}
      onKeyDown={e => onKeyDown(e)}
      onFocus={onChange}
      onBlur={onChange}
      {...props}
    />
  );
};

export default TextField;
