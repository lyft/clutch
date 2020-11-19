import * as React from "react";
import styled from "@emotion/styled";
import type { ButtonProps as MuiButtonProps, GridJustification } from "@material-ui/core";
import { Button as MuiButton, Grid, IconButton } from "@material-ui/core";
import CheckCircleOutlinedIcon from "@material-ui/icons/CheckCircleOutlined";
import FileCopyOutlinedIcon from "@material-ui/icons/FileCopyOutlined";

const COLORS = {
  neutral: {
    background: {
      primary: "#FFFFFF",
      hover: "#EFF0F2",
      active: "#CFD3D7",
      disabled: "#FFFFFF",
    },
    font: "#0D1030",
  },
  primary: {
    background: {
      primary: "#3548D4",
      hover: "#7680E8",
      active: "#1B31C6",
      disabled: "#7680E8",
    },
    font: "#FFFFFF",
  },
  danger: {
    background: {
      primary: "#DB3615",
      hover: "#FF6441",
      active: "#E32E11",
      disabled: "#FF3E1C",
    },
    font: "#FFFFFF",
  },
};

const StyledButton = styled(MuiButton)(
  {
    borderRadius: "4px",
    fontWeight: 700,
    textTransform: "none",
    padding: "0.875rem 2rem",
    margin: "2rem .875rem",
  },
  props => ({
    color: props["data-color"].font,
    backgroundColor: props["data-color"].background.primary,
    "&:hover": {
      backgroundColor: props["data-color"].background.hover,
    },
    "&:active": {
      backgroundColor: props["data-color"].background.active,
    },
    "&:disabled": {
      color: props["data-color"].font,
      backgroundColor: props["data-color"].background.disabled,
      opacity: "0.38",
    },
  })
);

const OutlinedButton = styled(StyledButton)({
  border: "1px solid #0D1030",
});

export interface ButtonProps
  extends Pick<MuiButtonProps, "disabled" | "endIcon" | "onClick" | "startIcon" | "type"> {
  // Case-sensitive button text.
  text: string;
  // Provides feedback to the user in regards to the action of the button.
  variant?: "neutral" | "primary" | "danger" | "destructive";
}

/**
 * A button with default themes based on use case.
 */
const Button: React.FC<ButtonProps> = ({ text, variant = "primary", ...props }) => {
  const color = variant === "destructive" ? "danger" : variant;

  const ButtonVariant = variant === "neutral" ? OutlinedButton : StyledButton;
  return (
    <ButtonVariant variant="contained" disableElevation data-color={COLORS[color]} {...props}>
      {text}
    </ButtonVariant>
  );
};

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
