import React from "react";
import styled from "@emotion/styled";
import { Fab, Link } from "@material-ui/core";
import ChatBubbleOutlineIcon from "@material-ui/icons/ChatBubbleOutline";

const Button = styled(Fab)({
  position: "fixed",
  right: "16px",
  bottom: "16px",
});

const Icon = styled(ChatBubbleOutlineIcon)({
  margin: "8px",
  color: "#2D3F50",
});

const FeedbackButton: React.FC<{}> = () => (
  // Hack because ButtonBaseProps can't take a component prop
  // https://github.com/mui-org/material-ui/issues/15827
  <Button variant="extended" size="small">
    <Link
      target="_blank"
      rel="noreferrer"
      href={process.env.REACT_APP_FEEDBACK_URL}
      underline="none"
      style={{ paddingTop: "25%" }}
    >
      <Icon fontSize="small" />
    </Link>
  </Button>
);

export default FeedbackButton;
