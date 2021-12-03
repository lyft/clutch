import React from "react";
import styled from "@emotion/styled";
import { Grow as MuiGrow, Paper as MuiPaper, Popper as MuiPopper } from "@material-ui/core";
import ChatBubbleOutlineIcon from "@material-ui/icons/ChatBubbleOutline";

import { IconButton } from "../button";

import NPSFeedback from "./feedback";

const Grow = styled(MuiGrow)((props: { placement: string }) => ({
  transformOrigin: props.placement,
}));

const Popper = styled(MuiPopper)({
  padding: "0 12px",
  marginLeft: "12px",
  zIndex: 1201,
});

const Paper = styled(MuiPaper)({
  width: "420px",
  border: "1px solid #E7E7EA",
  boxShadow: "0px 5px 15px rgba(53, 72, 212, 0.2)",
});

const StyledFeedbackIcon = styled(IconButton)<{ open: boolean }>(
  {
    color: "#ffffff",
    marginRight: "8px",
    padding: "12px",
    "&:hover": {
      background: "#2d3db4",
    },
    "&:active": {
      background: "#2938a5",
    },
  },
  props => ({
    background: props.open ? "#2d3db4" : "unset",
  })
);

const AnytimeFeedback = () => {
  const [open, setOpen] = React.useState<boolean>(false);
  const anchorRef = React.useRef(null);

  const handleToggle = () => {
    setOpen(!open);
  };

  const FeedbackIconProps = {
    edge: "end",
    id: "anytimeFeedbackIcon",
  };

  return (
    <>
      <StyledFeedbackIcon
        variant="neutral"
        ref={anchorRef}
        aria-controls={open ? "anytime-feedback" : undefined}
        open={open}
        aria-haspopup="true"
        onClick={handleToggle}
        {...FeedbackIconProps}
      >
        <ChatBubbleOutlineIcon />
      </StyledFeedbackIcon>
      <Popper open={open} anchorEl={anchorRef.current} transition placement="bottom-end">
        {({ TransitionProps, placement }) => (
          <Grow
            {...TransitionProps}
            placement={placement === "bottom" ? "center top" : "center bottom"}
          >
            <Paper>
              <NPSFeedback origin="ANYTIME" />
            </Paper>
          </Grow>
        )}
      </Popper>
    </>
  );
};

export default AnytimeFeedback;
