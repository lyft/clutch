import * as React from "react";
import type {
  ButtonProps as MuiButtonProps,
  GridJustification,
  IconButtonProps as MuiIconButtonProps,
} from "@material-ui/core";
import { Button as MuiButton, Grid, IconButton as MuiIconButton } from "@material-ui/core";
import CheckCircleOutlinedIcon from "@material-ui/icons/CheckCircleOutlined";
import FileCopyOutlinedIcon from "@material-ui/icons/FileCopyOutlined";

import { Tooltip } from "./Feedback/tooltip";
import styled from "./styled";

interface ButtonPalette {
  /** A palette of background colors used for the various button states. */
  background: {
    primary: string;
    hover: string;
    active: string;
    disabled: string;
  };
  /** A palette of font colors used for the various button states. */
  font: {
    primary: string;
    disabled?: string;
  };
}

const COLORS = {
  neutral: {
    background: {
      primary: "transparent",
      hover: "#E7E7EA",
      active: "#CFD3D7",
      disabled: "#FFFFFF",
    },
    font: {
      primary: "#0D1030",
      disabled: "#0D1030",
    },
  },
  primary: {
    background: {
      primary: "#3548D4",
      hover: "#2D3DB4",
      active: "#2938A5",
      disabled: "#E7E7EA",
    },
    font: {
      primary: "#FFFFFF",
      disabled: "rgba(13, 16, 48, 0.38)",
    },
  },
  danger: {
    background: {
      primary: "#DB3615",
      hover: "#BA2E12",
      active: "#AB2A10",
      disabled: "#F1B3A6",
    },
    font: {
      primary: "#FFFFFF",
      disabled: "#FFFFFF",
    },
  },
  secondary: {
    background: {
      primary: "transparent",
      hover: "#F5F6FD",
      active: "#D7DAF6",
      disabled: "transparent",
    },
    font: {
      primary: "#3548D4",
      disabled: "#0D1030",
    },
  },
} as { [key: string]: ButtonPalette };

const colorCss = (palette: ButtonPalette) => {
  return {
    color: palette.font.primary,
    backgroundColor: palette.background.primary,
    "&:hover": {
      backgroundColor: palette.background.hover,
    },
    "&:active": {
      backgroundColor: palette.background.active,
    },
    "&:disabled": {
      color: palette.font?.disabled ? palette.font?.disabled : palette.font?.primary,
      backgroundColor: palette.background.disabled,
      opacity: "0.38",
    },
  };
};

const StyledButton = styled(MuiButton)<{ palette: ButtonPalette }>(
  {
    borderRadius: "4px",
    fontWeight: 500,
    lineHeight: "20px",
    fontSize: "16px",
    textTransform: "none",
    height: "48px",
    padding: "14px 32px",
  },
  props => colorCss(props.palette)
);

const StyledBorderButton = styled(StyledButton)({
  border: "1px solid #0D1030",
  "&.Mui-disabled": {
    borderColor: "rgba(13, 16, 48, 0.1)",
  },
});

/** Provides feedback to the user in regards to the action of the button. */
type ButtonVariant = "neutral" | "primary" | "danger" | "destructive" | "secondary";

const ICON_BUTTON_STYLE_MAP = {
  small: {
    size: 32,
    padding: 4,
  },
  medium: {
    size: 48,
    padding: 12,
  },
  large: {
    size: 64,
    padding: 6,
  },
};
export type IconButtonSize = keyof typeof ICON_BUTTON_STYLE_MAP;

export const ICON_BUTTON_VARIANTS = Object.keys(ICON_BUTTON_STYLE_MAP);

/** A color palette from a @type ButtonPalette */
const variantPalette = (variant: ButtonVariant): ButtonPalette => {
  const v = variant === "destructive" ? "danger" : variant;
  return COLORS?.[v] || COLORS.primary;
};

export interface ButtonProps
  extends Pick<MuiButtonProps, "disabled" | "endIcon" | "onClick" | "startIcon" | "type"> {
  /** Case-sensitive button text. */
  text: string;
  /** The button variantion. Defaults to primary. */
  variant?: ButtonVariant;
}

/** A button with default themes based on use case. */
const Button = ({ text, variant = "primary", ...props }: ButtonProps) => {
  const palette = variantPalette(variant);
  const ButtonVariant = variant === "neutral" ? StyledBorderButton : StyledButton;

  return (
    <ButtonVariant variant="contained" disableElevation palette={palette} {...props}>
      {text}
    </ButtonVariant>
  );
};

const StyledIconButton = styled(MuiIconButton)<{
  $palette: ButtonPalette;
  $size?: IconButtonSize;
}>({}, props => ({
  width: `${ICON_BUTTON_STYLE_MAP[props.$size]?.size || ICON_BUTTON_STYLE_MAP.small.size}px`,
  height: `${ICON_BUTTON_STYLE_MAP[props.$size]?.size || ICON_BUTTON_STYLE_MAP.small.size}px`,
  padding: `${
    ICON_BUTTON_STYLE_MAP[props.$size]?.padding || ICON_BUTTON_STYLE_MAP.small.padding
  }px`,
  ...colorCss(props.$palette),
}));

// TODO: (jslaughter) Update when large sizing is available with material-ui@5
export interface IconButtonProps extends Pick<MuiIconButtonProps, "disabled" | "type" | "onClick"> {
  /** The button variantion. Defaults to primary. */
  variant?: ButtonVariant;
  children: React.ReactElement;
  size?: IconButtonSize;
}

/**
 * A button to wrap icons with default themes based on use case.
 * Will forwardRef so that tooltips can be wrapped around the buttons
 * @param variant valid color variant
 * @param size a valid size for the IconButton
 * @param children any children to render inside of the IconButton
 * @returns rendered IconButton component
 */
const IconButton = React.forwardRef<HTMLButtonElement, IconButtonProps>(
  ({ variant = "primary", size = "medium", children, ...props }: IconButtonProps, ref) => (
    <StyledIconButton $palette={variantPalette(variant)} $size={size} {...props} {...{ ref }}>
      {children}
    </StyledIconButton>
  )
);

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
  /** Buttons within the group. */
  children: React.ReactElement<ButtonProps> | React.ReactElement<ButtonProps>[];
  /** Justification of buttons. */
  justify?: GridJustification;
  /** Position of button group border. Defaults to top. */
  border?: "top" | "bottom";
}

/** A set of buttons to group together. */
const ButtonGroup = ({ children, justify = "flex-end", border = "top" }: ButtonGroupProps) => (
  <ButtonGroupContainer container justify={justify} data-border={border}>
    {children}
  </ButtonGroupContainer>
);

const StyledClipboardIconButton = styled(MuiIconButton)({
  color: "#000000",
  ":hover": {
    backgroundColor: "transparent",
  },
});

export interface ClipboardButtonProps {
  /** Case-sensitive text to be copied. */
  text: string;
  tooltip?: string;
}

/**
 * A button to copy text to a users clipboard.
 *
 * When clicked a checkmark is briefly displayed.
 */
const ClipboardButton = ({ text, tooltip = "" }: ClipboardButtonProps) => {
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
    <Tooltip title={tooltip}>
      <StyledClipboardIconButton
        onClick={() => {
          setClicked(true);
          navigator.clipboard.writeText(text);
        }}
      >
        {clicked ? <CheckCircleOutlinedIcon /> : <FileCopyOutlinedIcon />}
      </StyledClipboardIconButton>
    </Tooltip>
  );
};

export { Button, ButtonGroup, ClipboardButton, IconButton };
