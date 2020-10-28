import React from "react";
import type { ButtonProps as MuiButtonProps, GridJustification } from "@material-ui/core";
import { Button as MuiButton, emphasize, fade, Grid, IconButton } from "@material-ui/core";
import CheckCircleOutlinedIcon from "@material-ui/icons/CheckCircleOutlined";
import FileCopyOutlinedIcon from "@material-ui/icons/FileCopyOutlined";
import styled from "styled-components";

const StyledButton = styled(MuiButton)`
  ${({ theme, ...props }) => `
  margin: 15px;
  background-color: ${
    props["data-destructive"] === "true"
      ? theme.palette.destructive.main
      : fade(theme.palette.secondary.main, 0.3)
  };
  &:hover {
    background-color: ${emphasize(theme.palette.secondary.main, 0.3)};
  };
  `}
`;

interface ButtonProps
  extends Pick<
    MuiButtonProps,
    "disabled" | "endIcon" | "onClick" | "size" | "startIcon" | "type" | "variant"
  > {
  text: string;
  destructive?: boolean;
}

const Button: React.FC<ButtonProps> = ({ text, destructive = false, ...props }) => (
  <StyledButton variant="contained" data-destructive={destructive.toString()} {...props}>
    {text}
  </StyledButton>
);

const AdvanceButton: React.FC<ButtonProps> = ({ text, ...props }) => (
  <Button destructive={false} text={text} {...props} />
);
const DestructiveButton: React.FC<ButtonProps> = ({ text, ...props }) => (
  <Button destructive text={text} {...props} />
);

interface ButtonGroupProps {
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

interface ClipboardButtonProps {
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

export {
  AdvanceButton,
  Button,
  ButtonGroup,
  ButtonGroupProps,
  ButtonProps,
  ClipboardButton,
  DestructiveButton,
};
