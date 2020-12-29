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
      hover: "#E7E7EA",
      active: "#CFD3D7",
      disabled: "#FFFFFF",
    },
    font: "#0D1030",
  },
  primary: {
    background: {
      primary: "#3548D4",
      hover: "#2D3DB4",
      active: "#2938A5",
      disabled: "#E7E7EA",
    },
    font: "#FFFFFF",
  },
  danger: {
    background: {
      primary: "#DB3615",
      hover: "#BA2E12",
      active: "#AB2A10",
      disabled: "#F1B3A6",
    },
    font: "#FFFFFF",
  },
};

const StyledButton = styled(MuiButton)(
  {
    borderRadius: "4px",
    fontWeight: 500,
    lineHeight: "20px",
    fontSize: "16px",
    textTransform: "none",
    height: "48px",
    padding: "14px 32px",
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
  "&.Mui-disabled": {
    borderColor: "rgba(13, 16, 48, 0.1)",
  },
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

const ButtonGroupContainer = styled(Grid)(
  {
    "> *": {
      margin: "8px",
    },
  },
  props => (
    props["data-border"] === "top" ? {
    marginTop: "24px",
    borderTop: "1px solid #E7E7EA",
  } : {
    marginBottom: "24px",
    borderBottom: "1px solid #E7E7EA",
  })
);

export interface ButtonGroupProps {
  children: React.ReactElement<ButtonProps> | React.ReactElement<ButtonProps>[];
  justify?: GridJustification;
  border?: "top" | "bottom"
}

const ButtonGroup = ({ children, justify = "flex-end", border = "top" }: ButtonGroupProps) => (
  <ButtonGroupContainer container justify={justify} data-border={border}>
    {children}
  </ButtonGroupContainer>
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
