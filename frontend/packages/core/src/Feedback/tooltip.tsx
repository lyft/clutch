import * as React from "react";
import styled from "@emotion/styled";
import type { TooltipProps as MuiTooltipProps } from "@material-ui/core";
import { Tooltip as MuiTooltip } from "@material-ui/core";

const BaseTooltip = ({ className, ...props }: MuiTooltipProps) => (
  <MuiTooltip classes={{ tooltip: className }} {...props} />
);

const StyledTooltip = styled(BaseTooltip)(
  {
    backgroundColor: "#0D1030",
    border: "1px solid rgba(13, 16, 48, 0.1)",
    boxShadow: "0px 4px 6px rgba(53, 72, 212, 0.2)",
    borderRadius: "6px",
    margin: "2px",
  },
  props => ({
    maxWidth: props["data-max-width"],
  })
);

export interface TooltipProps extends Pick<MuiTooltipProps, "placement"> {
  // tooltip reference element (i.e. icon)
  children: React.ReactElement;
  // tooltip text
  title: React.ReactNode;
  // material ui default is 300px
  maxWidth?: string;
}

const Tooltip = ({ children, title, maxWidth = "300px", ...props }: TooltipProps) => {
  return (
    <StyledTooltip title={title} {...props} data-max-width={maxWidth}>
      {children}
    </StyledTooltip>
  );
};

// sets the spacing between multiline content
const TooltipContainer = styled.div({
  paddingBottom: "4px",
});

export { Tooltip, TooltipContainer };
