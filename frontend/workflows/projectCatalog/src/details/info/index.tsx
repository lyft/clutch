import React from "react";
import { Grid, styled } from "@clutch-sh/core";

import ChipsRow from "./chipsRow";
import LanguageRow from "./languageRow";
import MessengerRow from "./messengerRow";
import RepositoryRow from "./repositoryRow";
import type { ProjectInfo } from "./types";

interface ProjectInfoProps {
  data: ProjectInfo;
}

const StyledRow = styled(Grid)({
  marginBottom: "5px",
  whiteSpace: "nowrap",
  width: "100%",
});

const ProjectInfoCard = ({ data }: ProjectInfoProps) => (
  <>
    {data?.messenger && (
      <StyledRow container spacing={1}>
        <MessengerRow {...data.messenger} />
      </StyledRow>
    )}
    {data?.repository && (
      <StyledRow container spacing={1} justify="flex-start" alignItems="center">
        <RepositoryRow {...data.repository} />
      </StyledRow>
    )}
    {data?.languages?.length && (
      <StyledRow container spacing={1} justify="flex-start" alignItems="flex-end">
        <LanguageRow languages={data.languages} />
      </StyledRow>
    )}
    {data?.chips?.length && (
      <StyledRow container spacing={1}>
        <ChipsRow chips={data.chips} />
      </StyledRow>
    )}
  </>
);

export default ProjectInfoCard;
