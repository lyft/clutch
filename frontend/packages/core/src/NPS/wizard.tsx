import React, { useState } from "react";
import styled from "@emotion/styled";

import NPSFeedback from "./feedback";

const NPSContainer = styled.div<{ submit: boolean }>(
  {
    width: "50%",
    margin: "auto",
  },
  props => ({
    background: props.submit ? "unset" : "#F9F9FE",
  })
);

const NPSWizard = () => {
  const [hasSubmit, setSubmit] = useState<boolean>(false);

  const onSubmit = (submit: boolean) => {
    setSubmit(submit);
  };

  return (
    <NPSContainer submit={hasSubmit}>
      <NPSFeedback origin="WIZARD" onSubmit={onSubmit} />
    </NPSContainer>
  );
};

export default NPSWizard;
