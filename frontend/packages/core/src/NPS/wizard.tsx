import React, { useState } from "react";

import styled from "../styled";

import NPSFeedback from "./feedback";

const NPSContainer = styled("div")<{ $submit: boolean }>(
  {
    width: "50%",
    margin: "auto",
    borderRadius: "8px",
  },
  props => ({
    background: props.$submit ? "unset" : "#F9F9FE",
  })
);

const NPSWizard = () => {
  const [hasSubmit, setSubmit] = useState<boolean>(false);

  return (
    <NPSContainer $submit={hasSubmit} data-testid="nps-wizard">
      <NPSFeedback origin="WIZARD" onSubmit={setSubmit} />
    </NPSContainer>
  );
};

export default NPSWizard;
