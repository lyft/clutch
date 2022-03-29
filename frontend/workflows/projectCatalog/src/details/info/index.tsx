import React from "react";
import { Grid, styled } from "@clutch-sh/core";

import ChipsRow from "./chipsRow";
import LanguageRow from "./languageRow";
import MessengerRow from "./messengerRow";
import RepositoryRow from "./repositoryRow";
import type { ProjectInfo, ProjectInfoChip } from "./types";

interface ProjectInfoProps {
  data: ProjectInfo;
  addtlChips?: ProjectInfoChip[];
}

const StyledRow = styled(Grid)({
  marginBottom: "5px",
  whiteSpace: "nowrap",
  width: "100%",
});

const ProjectInfoCard = ({ data, addtlChips }: ProjectInfoProps) => {
  const [chips, setChips] = React.useState<ProjectInfoChip[]>([]);

  React.useEffect(() => {
    setChips((data?.chips || []).concat(addtlChips ?? []));
  }, [data, addtlChips]);

  return (
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
      {data?.languages?.length ? (
        <StyledRow container spacing={1} justify="flex-start" alignItems="flex-end">
          <LanguageRow languages={data.languages} />
        </StyledRow>
      ) : null}
      {chips.length > 0 && (
        <StyledRow container spacing={1}>
          <ChipsRow chips={chips} />
        </StyledRow>
      )}
    </>
  );
};

export default ProjectInfoCard;
