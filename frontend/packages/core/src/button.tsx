import * as React from "react";
import CheckCircleOutlinedIcon from "@mui/icons-material/CheckCircleOutlined";
import FileCopyOutlinedIcon from "@mui/icons-material/FileCopyOutlined";
import type {
  ButtonProps as MuiButtonProps,
  IconButtonProps as MuiIconButtonProps,
  Theme,
} from "@mui/material";
import {
  alpha,
  Button as MuiButton,
  Grid,
  IconButton as MuiIconButton,
  useTheme,
} from "@mui/material";

import { Tooltip } from "./Feedback/tooltip";
import type { GridJustification } from "./grid";
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

const BUTTON_SIZE_MAP = {
  xsmall: {
    height: "24px",
    padding: "4px 20px",
  },
  small: {
    height: "32px",
    padding: "7px 32px",
  },
  medium: {
    height: "48px",
    padding: "14px 32px",
  },
  large: {
    height: "64px",
    padding: "21px 32px",
  },
};

export type ButtonSize = keyof typeof BUTTON_SIZE_MAP;

const StyledButton = styled(MuiButton)<{
  palette: ButtonPalette;
  size: ButtonSize;
}>(
  {
    borderRadius: "4px",
    fontWeight: 500,
    lineHeight: "20px",
    fontSize: "16px",
    textTransform: "none",
    margin: "12px 8px",
    textOverflow: "ellipsis",
    whiteSpace: "nowrap",
    overflow: "hidden",
  },
  props => ({
    ...colorCss(props.palette),
    ...BUTTON_SIZE_MAP[props.size],
  })
);

const StyledBorderButton = styled(StyledButton)(({ theme }: { theme: Theme }) => ({
  border: `1px solid ${theme.palette.secondary[900]}`,
  "&.Mui-disabled": {
    borderColor: alpha(theme.palette.secondary[900], 0.1),
  },
}));

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
const variantPalette = (variant: ButtonVariant, theme: Theme): ButtonPalette => {
  const COLORS = {
    neutral: {
      background: {
        primary: "transparent",
        hover: theme.palette.secondary[200],
        active: theme.palette.secondary[300],
        disabled: theme.palette.contrastColor,
      },
      font: {
        primary: theme.palette.secondary[900],
        disabled: theme.palette.secondary[900],
      },
    },
    primary: {
      background: {
        primary: theme.palette.primary[600],
        hover: theme.palette.primary[700],
        active: theme.palette.primary[800],
        disabled: theme.palette.secondary[200],
      },
      font: {
        primary: theme.palette.contrastColor,
        disabled: alpha(theme.palette.secondary[900], 0.38),
      },
    },
    danger: {
      background: {
        primary: theme.palette.error[600],
        hover: theme.palette.error[700],
        active: theme.palette.error[800],
        disabled: theme.palette.error[200],
      },
      font: {
        primary: theme.palette.contrastColor,
        disabled: theme.palette.contrastColor,
      },
    },
    secondary: {
      background: {
        primary: "transparent",
        hover: theme.palette.primary[100],
        active: theme.palette.primary[300],
        disabled: "transparent",
      },
      font: {
        primary: theme.palette.primary[600],
        disabled: theme.palette.secondary[900],
      },
    },
  } as { [key: string]: ButtonPalette };
  const v = variant === "destructive" ? "danger" : variant;
  return COLORS?.[v] || COLORS.primary;
};

export interface ButtonProps
  extends Pick<MuiButtonProps, "id" | "disabled" | "endIcon" | "onClick" | "startIcon" | "type"> {
  /**
   * Case-sensitive button text
   */
  text: string;
  /**
   * The button variantion
   * @defaultValue primary
   */
  variant?: ButtonVariant;
  /**
   * The buttons size
   * @defaultValue medium
   */
  size?: ButtonSize;
}

/** A button with default themes based on use case. */
const Button = ({ text, variant = "primary", size = "medium", ...props }: ButtonProps) => {
  const theme = useTheme();
  const palette = variantPalette(variant, theme);
  const ButtonVariant = (variant === "neutral" ? StyledBorderButton : StyledButton) as any;

  return (
    <ButtonVariant variant="contained" disableElevation palette={palette} size={size} {...props}>
      {text}
    </ButtonVariant>
  );
};

const StyledIconButton = styled(MuiIconButton)<{
  $palette: ButtonPalette;
  $size?: MuiIconButtonProps["size"];
}>({}, props => ({
  width: `${ICON_BUTTON_STYLE_MAP[props.$size]?.size || ICON_BUTTON_STYLE_MAP.small.size}px`,
  height: `${ICON_BUTTON_STYLE_MAP[props.$size]?.size || ICON_BUTTON_STYLE_MAP.small.size}px`,
  padding: `${
    ICON_BUTTON_STYLE_MAP[props.$size]?.padding || ICON_BUTTON_STYLE_MAP.small.padding
  }px`,
  ...colorCss(props.$palette),
}));

export interface IconButtonProps
  extends Pick<MuiIconButtonProps, "disabled" | "type" | "onClick" | "size"> {
  /** The button variantion. Defaults to primary. */
  variant?: ButtonVariant;
  children: React.ReactElement;
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
  ({ variant = "primary", size = "medium", children, ...props }: IconButtonProps, ref) => {
    const theme = useTheme();

    return (
      <StyledIconButton
        $palette={variantPalette(variant, theme)}
        $size={size}
        {...props}
        {...{ ref }}
      >
        {children}
      </StyledIconButton>
    );
  }
);

const ButtonGroupContainer = styled(Grid)(
  {
    "> *": {
      margin: "12px 8px",
    },
  },
  props => ({ theme }: { theme: Theme }) =>
    props["data-border"] === "bottom"
      ? {
          marginBottom: "12px",
          borderBottom: `1px solid ${theme.palette.secondary[200]}`,
          marginTop: "0",
        }
      : {
          marginTop: "12px",
          borderTop: `1px solid ${theme.palette.secondary[200]}`,
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
  <ButtonGroupContainer container justifyContent={justify} data-border={border}>
    {children}
  </ButtonGroupContainer>
);

const StyledClipboardIconButton = styled(MuiIconButton)(({ theme }: { theme: Theme }) => ({
  color: theme.palette.getContrastText(theme.palette.contrastColor),
  ":hover": {
    backgroundColor: "transparent",
  },
}));

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
