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
              are user-owned or searched services, libraries, mobile apps, etc. Unchecking a project hides its respective upstream
              and downstream dependencies.
            </StyledTooltipBody>
          </StyledToolTipContainer>
          <StyledToolTipContainer>
            <StyledTooltipTitle>Upstreams - </StyledTooltipTitle>
            <StyledTooltipBody>
              are dependencies that are relied on by the respective projects (i.e. a service calls into this upstream for data).
            </StyledTooltipBody>
          </StyledToolTipContainer>
          <StyledToolTipContainer>
            <StyledTooltipTitle>Downstreams - </StyledTooltipTitle>
            <StyledTooltipBody>
              are dependencies that rely on the respective projects (i.e. this downstream calls into a service for data).
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
