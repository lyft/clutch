import React from "react";
import { ClutchError, Grid } from "@clutch-sh/core";

import DeployEventIcon from "../../assets/DeployEvent";
import type { BaseProjectCardProps } from "../card";
import ProjectCard, { LastEvent, StyledLink } from "../card";

import CommitInformation from "./commitInformation";
import type { CommitInfo, ProjectDeploys } from "./types";

interface ProjectDeploysProps {
  data: ProjectDeploys;
  error?: ClutchError | undefined;
  loading?: boolean;
}

const Deploys = ({ deploys }: { deploys: CommitInfo[] }) => (
  <>
    {deploys.map(deploy => (
      <Grid item xs={12}>
        <CommitInformation {...deploy} />
      </Grid>
    ))}
  </>
);

const ProjectDeploysCard = ({ data, error, loading }: ProjectDeploysProps) => {
  const titleData: BaseProjectCardProps = {
    text: data?.title ?? "Deploys",
    icon: <DeployEventIcon />,
    endAdornment: <LastEvent time={data?.lastDeploy} />,
  };

  return (
    <ProjectCard loading={loading} error={error} {...titleData}>
      {data?.deploys?.length && (
        <Grid container direction="row" spacing={2}>
          <Deploys deploys={data?.deploys} />
        </Grid>
      )}
      {data?.seeMore && (
        <Grid container item direction="column" alignItems="flex-end" style={{ marginTop: "10px" }}>
          <Grid item xs={6}>
            <StyledLink href={data.seeMore.url}>{data.seeMore.text}</StyledLink>
          </Grid>
        </Grid>
      )}
    </ProjectCard>
  );
};

export default ProjectDeploysCard;
