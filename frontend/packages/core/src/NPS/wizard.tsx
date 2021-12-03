import React from "react";
import styled from "@emotion/styled";

import NPSFeedback from "./feedback";

const NPSContainer = styled.div({
  width: "50%",
  margin: "auto",
  background: "#F9F9FE",
});

const NPSWizard = () => (
  <NPSContainer>
    <NPSFeedback origin="WIZARD" />
  </NPSContainer>
);

export default NPSWizard;
