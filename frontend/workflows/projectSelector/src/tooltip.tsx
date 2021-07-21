import * as React from "react";
import styled from "@emotion/styled";
import { Popper as MuiPopper, Tooltip } from "@material-ui/core";
import InfoOutlinedIcon from "@material-ui/icons/InfoOutlined";

const Popper = styled(MuiPopper)({
  ".MuiTooltip-tooltip": {
    color: "rgba(13, 16, 48, 0.6)",
    backgroundColor: "#FFFFFF",
    border: "1px solid rgba(13, 16, 48, 0.1)",
    boxShadow: "0px 4px 6px rgba(53, 72, 212, 0.2)",
    borderRadius: "6px",
    minWidth: "554px",
    maxWidth: "554px",
    padding: "0px",
  },
  ".MuiTooltip-tooltipPlacementRight": {
    margin: "0 4px",
  },
});

const renderPopper = props => {
  return <Popper {...props} />;
};

const StyledToolTipContainer = styled.div({
  padding: "8px 16px 8px 16px",
});

const StyledTooltipTitle = styled.span({
  fontWeight: 500,
  fontSize: "14px",
  lineHeight: "18px",
});

const StyledTooltipBody = styled.span({
  fontWeight: 400,
  fontSize: "12px",
  lineHeight: "16px",
});

const StyledTooltip = () => {
  return (
    <Tooltip
      title={
        <>
          <StyledToolTipContainer>
            <StyledTooltipTitle>Projects - </StyledTooltipTitle>
            <StyledTooltipBody>
              Lorem ipsum dolor sit amet, consectetur adipiscing elit. Pellentesque laoreet
              tristique pharetra, eu.
            </StyledTooltipBody>
          </StyledToolTipContainer>
          <StyledToolTipContainer>
            <StyledTooltipTitle>Upstreams - </StyledTooltipTitle>
            <StyledTooltipBody>
              Lorem ipsum dolor sit amet, consectetur adipiscing elit. Pellentesque laoreet
              tristique pharetra, eu.
            </StyledTooltipBody>
          </StyledToolTipContainer>
          <StyledToolTipContainer>
            <StyledTooltipTitle>Downstreams - </StyledTooltipTitle>
            <StyledTooltipBody>
              Lorem ipsum dolor sit amet, consectetur adipiscing elit. Pellentesque laoreet
              tristique pharetra, eu.
            </StyledTooltipBody>
          </StyledToolTipContainer>
        </>
      }
      placement="right-start"
      PopperComponent={renderPopper}
    >
      <InfoOutlinedIcon />
    </Tooltip>
  );
};

export default StyledTooltip;
