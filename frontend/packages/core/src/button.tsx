import * as React from "react";
import styled from "@emotion/styled";
import type {
  ButtonProps as MuiButtonProps,
  GridJustification,
  IconButtonProps as MuiIconButtonProps,
} from "@material-ui/core";
import { Button as MuiButton, Grid, IconButton as MuiIconButton } from "@material-ui/core";
import CheckCircleOutlinedIcon from "@material-ui/icons/CheckCircleOutlined";
import FileCopyOutlinedIcon from "@material-ui/icons/FileCopyOutlined";

interface ButtonColor {
  background: {
    primary: string;
    hover: string;
    active: string;
    disabled: string;
  };
  font: string;
  fontDisabled: string;
}

const COLORS = {
  neutral: {
    background: {
      primary: "transparent",
      hover: "#E7E7EA",
      active: "#CFD3D7",
      disabled: "#FFFFFF",
    },
    font: "#0D1030",
    fontDisabled: "#0D1030",
  },
  primary: {
    background: {
      primary: "#3548D4",
      hover: "#2D3DB4",
      active: "#2938A5",
      disabled: "#E7E7EA",
    },
    font: "#FFFFFF",
    fontDisabled: "rgba(13, 16, 48, 0.38)",
  },
  danger: {
    background: {
      primary: "#DB3615",
      hover: "#BA2E12",
      active: "#AB2A10",
      disabled: "#F1B3A6",
    },
    font: "#FFFFFF",
    fontDisabled: "#FFFFFF",
  },
} as { [key: string]: ButtonColor };

const StyledButton = styled(MuiButton)<{ palette: ButtonColor }>(
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
    color: props.palette.font,
    backgroundColor: props.palette.background.primary,
    "&:hover": {
      backgroundColor: props.palette.background.hover,
    },
    "&:active": {
      backgroundColor: props.palette.background.active,
    },
    "&:disabled": {
      color: props.palette.fontDisabled,
      backgroundColor: props.palette.background.disabled,
      opacity: "0.38",
    },
  })
);

const StyledBorderButton = styled(StyledButton)({
  border: "1px solid #0D1030",
  "&.Mui-disabled": {
    borderColor: "rgba(13, 16, 48, 0.1)",
  },
});

type ButtonVariant = "neutral" | "primary" | "danger" | "destructive";

export interface ButtonProps
  extends Pick<MuiButtonProps, "disabled" | "endIcon" | "onClick" | "startIcon" | "type"> {
  /* Case-sensitive button text. */
  text: string;
  /* Provides feedback to the user in regards to the action of the button. */
  variant?: ButtonVariant;
}

/**
 * A button with default themes based on use case.
 */
const Button: React.FC<ButtonProps> = ({ text, variant = "primary", ...props }) => {
  const color = variant === "destructive" ? "danger" : variant;

  const ButtonVariant = variant === "neutral" ? StyledBorderButton : StyledButton;
  return (
    <ButtonVariant variant="contained" disableElevation palette={COLORS[color]} {...props}>
      {text}
    </ButtonVariant>
  );
};

const StyledIconButton = styled(MuiIconButton)<{ palette: ButtonColor }>(
  {
    height: "48px",
    width: "48px",
    padding: "12px",
  },
  props => ({
    color: props.palette.font,
    backgroundColor: props.palette.background.primary,
    "&:hover": {
      backgroundColor: props.palette.background.hover,
    },
    "&:active": {
      backgroundColor: props.palette.background.active,
    },
    "&:disabled": {
      color: "rgba(13, 16, 48, 0.38)",
      backgroundColor: props.palette.background.disabled,
      opacity: "0.38",
    },
  })
);

export interface IconButtonProps extends Pick<MuiIconButtonProps, "disabled" | "type" | "onClick"> {
  /* Provides feedback to the user in regards to the action of the button. */
  variant?: ButtonVariant;
  children: React.ReactElement;
}

const IconButton = ({ variant = "primary", children, ...props }: IconButtonProps) => {
  const color = variant === "destructive" ? "danger" : variant;

  return (
    <StyledIconButton {...props} palette={COLORS[color]}>
      {children}
    </StyledIconButton>
  );
};

const ButtonGroupContainer = styled(Grid)(
  {
    "> *": {
      margin: "12px 8px",
    },
  },
  props =>
    props["data-border"] === "bottom"
      ? {
          marginBottom: "12px",
          borderBottom: "1px solid #E7E7EA",
          marginTop: "0",
        }
      : {
          marginTop: "12px",
          borderTop: "1px solid #E7E7EA",
          marginBottom: "0",
        }
);

export interface ButtonGroupProps {
  children: React.ReactElement<ButtonProps> | React.ReactElement<ButtonProps>[];
  justify?: GridJustification;
  border?: "top" | "bottom";
}

const ButtonGroup = ({ children, justify = "flex-end", border = "top" }: ButtonGroupProps) => (
  <ButtonGroupContainer container justify={justify} data-border={border}>
    {children}
  </ButtonGroupContainer>
);

const ClipboardIconButton = styled(MuiIconButton)({
  color: "#000000",
  ":hover": {
    backgroundColor: "transparent",
  },
});

export interface ClipboardButtonProps {
  text: string;
}

// ClipboardButton is a button that copies text to the clipboard and briefly displays a checkmark
// after being clicked to let the user know that clicking actually did something and sent the
// provided value to the clipboard.
const ClipboardButton: React.FC<ClipboardButtonProps> = ({ text }) => {
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
    <ClipboardIconButton
      onClick={() => {
        setClicked(true);
        navigator.clipboard.writeText(text);
      }}
    >
      {clicked ? <CheckCircleOutlinedIcon /> : <FileCopyOutlinedIcon />}
    </ClipboardIconButton>
  );
};

export { Button, ButtonGroup, ClipboardButton, IconButton };
