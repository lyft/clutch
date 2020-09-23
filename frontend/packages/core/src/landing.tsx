import React from "react";
import { useNavigate } from "react-router-dom";
import { Card, CardActionArea, CardContent, Grid, Link, Paper } from "@material-ui/core";
import Typography from "@material-ui/core/Typography";
import GitHubIcon from "@material-ui/icons/GitHub";
import styled from "styled-components";

import { userId } from "./AppLayout/user";
import { useAppContext } from "./Contexts";
import { TrendingUpIcon } from "./icon";

const GridContainer = styled(Grid)`
  margin-top: 20px;
`;

const SizedCard = styled(CardContent)`
  ${({ theme }) => `
  height: 150px;
  width: 300px;
  background: linear-gradient(340deg, ${theme.palette.accent.main} 0%, ${theme.palette.secondary.main} 90%);
  `};
`;

const CardTitle = styled(Typography)`
  color: #ffffff;
`;

interface MediaCardProps {
  title: string;
  description: string;
  path: string;
}

const MediaCard: React.FC<MediaCardProps> = ({ title, description, path }) => {
  const navigate = useNavigate();
  const navigateTo = () => {
    navigate(path);
  };

  return (
    <Grid item>
      <Card raised>
        <CardActionArea onClick={navigateTo}>
          <SizedCard>
            <CardTitle gutterBottom variant="h6">
              {title}
            </CardTitle>
            <Typography variant="body2" color="textSecondary">
              {description}
            </Typography>
          </SizedCard>
        </CardActionArea>
      </Card>
    </Grid>
  );
};

const Content = styled(Paper)`
  padding: 1.5%;
`;

const Footer = styled.div`
  @media screen and (min-width: 900px) and (min-height: 500px) {
    position: absolute;
    bottom: 0;
    width: 150px;
    left: 50%;
    margin-left: -75px;
    padding-bottom: 10px;
  }
  @media screen and (max-height: 500px) {
    position: inherit;
  }
`;

const GitHubLogo = styled(GitHubIcon)`
  ${({ theme }) => `
  color: ${theme.palette.accent.main};
  margin-right: 5px;
  `}
`;

const GitHubLink = styled(Link)`
  ${({ theme }) => `
  color: ${theme.palette.secondary.main};
  `}
`;

const Landing: React.FC<{}> = () => {
  const { workflows } = useAppContext();
  const trendingWorkflows = [];
  workflows.forEach(workflow => {
    workflow.routes.forEach(route => {
      const title = route.displayName
        ? `${workflow.displayName} - ${route.displayName}`
        : workflow.displayName;
      if (route.trending) {
        trendingWorkflows.push({
          title,
          description: route.description,
          path: `${workflow.path}/${route.path}`,
        });
      }
    });
  });

  return (
    <Content id="landing" elevation={0}>
      <Typography variant="h5">
        <strong>Welcome {userId()} </strong>
        <span role="img" aria-label="Hand Waving">
          👋
        </span>
      </Typography>
      <>
        <div>
          <Typography gutterBottom variant="body1" paragraph>
            Clutch will assist you in safely modifying resources outside of the normal orchestration
            process.
          </Typography>

          {trendingWorkflows.length === 0 ? null : (
            <>
              <Grid container justify="center" alignItems="center">
                <TrendingUpIcon />
                <Typography align="center" variant="h5">
                  Trending Workflows
                </Typography>
              </Grid>

              <GridContainer justify="center" container direction="row" spacing={3}>
                {trendingWorkflows.map(workflow => (
                  <MediaCard
                    title={workflow.title}
                    description={workflow.description}
                    path={workflow.path}
                    key={workflow.path}
                  />
                ))}
              </GridContainer>
            </>
          )}
        </div>
        <Footer>
          <GridContainer container justify="center">
            <Grid item>
              <GitHubLogo fontSize="small" />
            </Grid>
            <Grid item>
              <GitHubLink
                target="_blank"
                rel="noreferrer"
                href="https://github.com/lyft/clutch"
                underline="none"
              >
                lyft/clutch
              </GitHubLink>
            </Grid>
          </GridContainer>
        </Footer>
      </>
    </Content>
  );
};

export default Landing;
