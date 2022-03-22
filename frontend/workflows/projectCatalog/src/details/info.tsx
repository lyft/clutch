import React from "react";
import { Card, Chip, Grid, Link, styled, Tooltip, Typography } from "@clutch-sh/core";
import { faGithub, faSlack } from "@fortawesome/free-brands-svg-icons";
import { faCodeBranch, faLock } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import GroupIcon from "@material-ui/icons/Group";
import SecurityIcon from "@material-ui/icons/Security";

import LanguageIcon from "../helpers/language-icon";

interface ProjectInfoProps {
  owners?: string[];
  repository?: string;
  slo?: number;
  slack?: string;
  languages?: string[];
  disabled?: boolean;
  name: string;
  tier: string;
}

const StyledCard = styled(Card)({
  width: "100%",
  height: "fit-content",
  padding: "15px",
});

const StyledRow = styled(Grid)({
  marginBottom: "5px",
  whiteSpace: "nowrap",
  width: "100%",
});

const LinkText = ({ text, link }: { text: string; link?: string }) => {
  const returnText = <Typography variant="body2">{text}</Typography>;

  if (link && text) {
    return <Link href={link}>{returnText}</Link>;
  }

  return returnText;
};

const ProjectInfo = ({
  owners = [],
  slack,
  name,
  repository,
  languages,
  tier,
  slo,
  disabled,
}: ProjectInfoProps) => {
  const windshieldScore = 94;
  const sloScore = 79;
  const repo = repository.split("/")[1].split(".git")[0];
  const pullRequests = 0;

  return (
    <StyledCard container direction="column">
      <StyledRow container spacing={1} justify="flex-start" alignItems="flex-start">
        <Grid item xs={1}>
          <GroupIcon />
        </Grid>
        <Grid item xs={10}>
          <Typography variant="subtitle1">{owners[0]}</Typography>
        </Grid>
        {disabled && (
          <Grid item xs={1}>
            <Tooltip title={`${name} is disabled`}>
              <FontAwesomeIcon icon={faLock} size="lg" />
            </Tooltip>
          </Grid>
        )}
      </StyledRow>
      <StyledRow container spacing={1}>
        {/* This should be more generic */}
        <Grid item>
          <FontAwesomeIcon icon={faSlack} size="lg" />
        </Grid>
        <Grid item>
          <LinkText text={slack} />
        </Grid>
      </StyledRow>
      <StyledRow container spacing={1} justify="flex-start" alignItems="center">
        <Grid item>
          <FontAwesomeIcon icon={faGithub} size="lg" />
        </Grid>
        <Grid item>
          <LinkText text={repo} />
        </Grid>
        {pullRequests > 0 && (
          <>
            <Grid item>
              <FontAwesomeIcon icon={faCodeBranch} size="1x" />
            </Grid>
            <Grid item>
              {/* Color Red */}
              <LinkText text={`${pullRequests} Open`} />
            </Grid>
          </>
        )}
      </StyledRow>
      {languages?.length && (
        <StyledRow container spacing={1} justify="flex-start" alignItems="flex-end">
          <Grid item>
            <Typography variant="body2">Languages</Typography>
          </Grid>
          {languages.map(language => (
            <Grid item>
              <LanguageIcon language={language} />
            </Grid>
          ))}
        </StyledRow>
      )}
      <Grid container spacing={1}>
        <Grid item>
          <Tooltip title={`Tier ${tier} Service`}>
            <Chip variant="neutral" label={`T${tier}`} size="small" />
          </Tooltip>
        </Grid>
        <Grid item>
          {/* Icon? How to make this extensible */}
          <Tooltip title={`SLO Score ${sloScore}%`}>
            <Chip variant="warn" label={`SLO ${sloScore}%`} size="small" />
          </Tooltip>
        </Grid>
        <Grid item>
          {/* Icon? */}
          <Tooltip title={`Windshield score ${windshieldScore}%`}>
            <Chip
              variant="warn"
              label={`${windshieldScore}%`}
              size="small"
              icon={<SecurityIcon />}
            />
          </Tooltip>
        </Grid>
      </Grid>
    </StyledCard>
  );
};

export default ProjectInfo;
