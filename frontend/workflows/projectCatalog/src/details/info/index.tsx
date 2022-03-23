import React from "react";
import type { ClutchError } from "@clutch-sh/core";
import { Grid, styled, Tooltip } from "@clutch-sh/core";
import { faLock } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import GroupIcon from "@material-ui/icons/Group";
import { capitalize } from "lodash";

import type { BaseProjectCardProps } from "../card";
import ProjectCard from "../card";

import ChipsRow from "./chipsRow";
import LanguageRow from "./languageRow";
import MessengerRow from "./messengerRow";
import RepositoryRow from "./repositoryRow";
import type { ProjectInfo } from "./types";

interface ProjectInfoProps {
  data: ProjectInfo;
  error?: ClutchError | undefined;
  loading?: boolean;
}

const StyledRow = styled(Grid)({
  marginBottom: "5px",
  whiteSpace: "nowrap",
  width: "100%",
});

const DisabledItem = ({ name }: { name: string }) => (
  <Grid item>
    <Tooltip title={`${name} is disabled`}>
      <FontAwesomeIcon icon={faLock} size="lg" />
    </Tooltip>
  </Grid>
);

const ProjectInfoCard = ({ data, loading, error }: ProjectInfoProps) => {
  const [titleData, setTitleData] = React.useState<BaseProjectCardProps>(null);

  React.useEffect(() => {
    if (data) {
      const capitalized = capitalize(data.name);
      setTitleData({
        text: capitalized,
        icon: <GroupIcon />,
        endAdornment: data.disabled ? <DisabledItem name={capitalized} /> : null,
      });
    }
  }, [data]);

  return (
    <>
      {titleData && (
        <ProjectCard loading={loading} error={error} {...titleData}>
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
        </ProjectCard>
      )}
    </>
  );
};

export default ProjectInfoCard;
