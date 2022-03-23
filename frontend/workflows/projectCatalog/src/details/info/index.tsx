import React from "react";
import type { ClutchError } from "@clutch-sh/core";
import { Grid, styled, Tooltip } from "@clutch-sh/core";
import { faLock } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import GroupIcon from "@material-ui/icons/Group";
import { capitalize } from "lodash";

import ProjectCard, { TitleRowProps } from "../card";

import ChipsRow from "./chipsRow";
import LanguageRow from "./languageRow";
import MessengerRow from "./messengerRow";
import RepositoryRow from "./repositoryRow";
import type { ProjectInfo } from "./types";

interface ProjectInfoProps {
  info: ProjectInfo;
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

const ProjectInfoCard = ({ info, loading }: ProjectInfoProps) => {
  const [titleData, setTitleData] = React.useState<TitleRowProps>(null);

  React.useEffect(() => {
    if (info) {
      const capitalized = capitalize(info.name);
      setTitleData({
        text: capitalized,
        icon: <GroupIcon />,
        endAdornment: info.disabled ? <DisabledItem name={capitalized} /> : null,
      });
    }
  }, [info]);

  return (
    <>
      {titleData && (
        <ProjectCard {...titleData}>
          {info?.messenger && (
            <StyledRow container spacing={1}>
              <MessengerRow {...info.messenger} />
            </StyledRow>
          )}
          {info?.repository && (
            <StyledRow container spacing={1} justify="flex-start" alignItems="center">
              <RepositoryRow {...info.repository} />
            </StyledRow>
          )}
          {info?.languages?.length && (
            <StyledRow container spacing={1} justify="flex-start" alignItems="flex-end">
              <LanguageRow languages={info.languages} />
            </StyledRow>
          )}
          {info?.chips?.length && (
            <StyledRow container spacing={1}>
              <ChipsRow chips={info.chips} />
            </StyledRow>
          )}
        </ProjectCard>
      )}
    </>
  );
};

export default ProjectInfoCard;
