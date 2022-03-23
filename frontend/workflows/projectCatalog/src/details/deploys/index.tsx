import React from "react";
import { Grid } from "@clutch-sh/core";

import DeployEventIcon from "../../assets/DeployEvent";
import ProjectCard, { LastEvent, StyledLink } from "../card";

import type { CommitInfo } from "./commitInformation";
import CommitInformation from "./commitInformation";

const Deploys = ({ deploys }: { deploys: CommitInfo[] }) => (
  <>
    {deploys.map(deploy => (
      <Grid item xs={12}>
        <CommitInformation {...deploy} />
      </Grid>
    ))}
  </>
);

export interface ProjectDeploys {
  title?: string;
  lastDeploy?: number;
  deploys: CommitInfo[];
  seeMore: {
    text: string;
    url: string;
  };
}

const ProjectDeploysCard = ({
  title = "Deploys",
  lastDeploy,
  seeMore,
  deploys,
}: ProjectDeploys) => {
  const titleData = {
    text: title,
    icon: <DeployEventIcon />,
    endAdornment: <LastEvent time={lastDeploy} />,
  };

  return (
    <ProjectCard {...titleData}>
      {deploys?.length && (
        <Grid container direction="row" spacing={2}>
          <Deploys deploys={deploys} />
        </Grid>
      )}
      {seeMore && (
        <Grid container item direction="column" alignItems="flex-end" style={{ marginTop: "10px" }}>
          <Grid item xs={6}>
            <StyledLink href={seeMore.url}>{seeMore.text}</StyledLink>
          </Grid>
        </Grid>
      )}
    </ProjectCard>
  );
};

export default ProjectDeploysCard;
