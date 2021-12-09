import * as React from "react";
import styled from "@emotion/styled";
import type { TooltipProps as MuiTooltipProps } from "@material-ui/core";
import { Tooltip as MuiTooltip } from "@material-ui/core";

const BaseTooltip = ({ className, ...props }: MuiTooltipProps) => (
  <MuiTooltip classes={{ tooltip: className }} {...props} />
);

// TODO: sync with Design on margins for each possible placement
const StyledTooltip = styled(BaseTooltip)((props: { maxwidth?: string }) => ({
  maxWidth: props.maxwidth,
  backgroundColor: "#0D1030",
  borderRadius: "6px",
  "&.MuiTooltip-tooltipPlacementLeft": {
    margin: "0 2px",
  },
  "&.MuiTooltip-tooltipPlacementRight": {
    margin: "0 2px",
  },
  "&.MuiTooltip-tooltipPlacementTop": {
    margin: "2px 0",
  },
  "&.MuiTooltip-tooltipPlacementBottom": {
    margin: "2px 0",
  },
}));

export interface TooltipProps extends Pick<MuiTooltipProps, "interactive" | "placement"> {
  // tooltip reference element (i.e. icon)
  children: React.ReactElement;
  // material ui default is 300px
  maxWidth?: string;
  // tooltip text
  title: React.ReactNode;
}

const Tooltip = ({ children, maxWidth = "300px", title, ...props }: TooltipProps) => {
  return (
    <StyledTooltip title={title} maxwidth={maxWidth} {...props}>
      {children}
    </StyledTooltip>
  );
};

// sets the spacing between multiline content
const TooltipContainer = styled.div({
  margin: "4px 0px",
});

export { Tooltip, TooltipContainer };
