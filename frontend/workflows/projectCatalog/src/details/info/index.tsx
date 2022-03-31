import React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";
import { Grid, styled } from "@clutch-sh/core";

import type { ProjectInfoChip } from "./chipsRow";
import ChipsRow from "./chipsRow";
import LanguageRow from "./languageRow";
import MessengerRow from "./messengerRow";
import RepositoryRow from "./repositoryRow";

interface ProjectInfoProps {
  projectData: IClutch.core.project.v1.IProject;
  addtlChips?: ProjectInfoChip[];
}

const StyledRow = styled(Grid)({
  marginBottom: "8px",
  whiteSpace: "nowrap",
  width: "100%",
});

const ProjectInfoCard = ({ projectData, addtlChips }: ProjectInfoProps) => {
  const [chips, setChips] = React.useState<ProjectInfoChip[]>([]);

  React.useEffect(() => {
    let tempChips: ProjectInfoChip[] = [];

    const { tier } = projectData;
    if (tier) {
      tempChips.push({
        text: `T${tier}`,
        title: `Tier ${tier} Service`,
      });
    }

    if (addtlChips) {
      tempChips = tempChips.concat(addtlChips);
    }

    setChips(tempChips);
  }, [projectData, addtlChips]);

  return (
    <>
      {projectData?.data && (
        <StyledRow container spacing={2}>
          <MessengerRow projectData={projectData} />
        </StyledRow>
      )}
      {projectData?.data?.repository && (
        <StyledRow container spacing={2} justify="flex-start" alignItems="center">
          <RepositoryRow repo={projectData.data.repository as string} />
        </StyledRow>
      )}
      {projectData?.languages?.length ? (
        <StyledRow container spacing={2} justify="flex-start" alignItems="flex-end">
          <LanguageRow languages={projectData.languages} />
        </StyledRow>
      ) : null}
      {chips.length ? (
        <Grid container spacing={2}>
          <ChipsRow chips={chips} />
        </Grid>
      ) : null}
    </>
  );
};

export default ProjectInfoCard;
