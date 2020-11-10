import React from "react";
import type { ButtonProps as MuiButtonProps, GridJustification } from "@material-ui/core";
import { Button as MuiButton, emphasize, fade, Grid, IconButton } from "@material-ui/core";
import CheckCircleOutlinedIcon from "@material-ui/icons/CheckCircleOutlined";
import FileCopyOutlinedIcon from "@material-ui/icons/FileCopyOutlined";
import styled from "styled-components";

const COLORS = {
  neutral: {
    background: "#FFFFFF",
    font: "#0D1030",
  },
  primary: {
    background: "#3B73E0",
    font: "#FFFFFF",
  },
  caution: {
    background: "#DA1707",
    font: "#FFFFFF",
  },
};

const StyledButton = styled(MuiButton)`
  ${({ theme, ...props }) => `
  margin: 15px;
  background-color: ${props["data-color"].background};
  color: ${props["data-color"].font};
  font-weight: 600;
  text-transform: none;
  height: 3rem;
  width: 8.875rem;
  margin: 2rem .875rem;
  &:hover {
    background-color: ${emphasize(props["data-color"].background, 0.2)};
  };
  &:disabled {
    color: ${props["data-color"].font};
    background-color: ${emphasize(props["data-color"].background, 0.6)};
  };
  `}
`;

export interface ButtonProps
  extends Pick<
    MuiButtonProps,
    "disabled" | "endIcon" | "onClick" | "size" | "startIcon" | "type"
  > {
  text: string;
  variant?: "neutral" | "primary" | "caution";
}

const Button: React.FC<ButtonProps> = ({ text, variant = "primary", ...props }) => (
  <StyledButton variant="contained" data-color={COLORS[variant]} {...props}>
    {text}
  </StyledButton>
);

export interface ButtonGroupProps {
  buttons: ButtonProps[];
  justify?: GridJustification;
}

const ButtonGroup: React.FC<ButtonGroupProps> = ({ buttons, justify = "center" }) => (
  <Grid container justify={justify}>
    {buttons.map(button => (
      <Button key={button.text} {...button} />
    ))}
  </Grid>
);

export interface ClipboardButtonProps {
  primary?: boolean;
  size?: "small" | "medium";
  text: string;
}

// ClipboardButton is a button that copies text to the clipboard and briefly displays a checkmark
// after being clicked to let the user know that clicking actually did something and sent the
// provided value to the clipboard.
const ClipboardButton: React.FC<ClipboardButtonProps> = ({
  text,
  primary = false,
  size = "small",
  ...props
}) => {
  const [clicked, setClicked] = React.useState(false);
  React.useEffect(() => {
    if (clicked) {
      const ticker = setTimeout(() => {
        setClicked(false);
      }, 1000);
      return () => clearTimeout(ticker);
    }

    return () => {};
  }, [clicked]);

  return (
    <IconButton
      color={primary ? "primary" : "secondary"}
      onClick={() => {
        setClicked(true);
        navigator.clipboard.writeText(text);
      }}
      size={size}
      {...props}
    >
      {clicked ? <CheckCircleOutlinedIcon /> : <FileCopyOutlinedIcon />}
    </IconButton>
  );
};

export { Button, ButtonGroup, ClipboardButton };
