import React from "react";
import { Fab, Link } from "@material-ui/core";
import ChatBubbleOutlineIcon from "@material-ui/icons/ChatBubbleOutline";
import styled from "styled-components";

const Button = styled(Fab)`
  ${({ theme }) => `
  position: fixed;
  right: ${theme.spacing(2)}px;
  bottom: ${theme.spacing(2)}px;
  `}
`;

const Icon = styled(ChatBubbleOutlineIcon)`
  ${({ theme }) => `
  margin: ${theme.spacing(1)}px;
  color: ${theme.palette.secondary.main}
  `}
`;

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
