import React from "react";
import { useNavigate } from "react-router-dom";
import styled from "@emotion/styled";
import { Grid, Paper } from "@material-ui/core";
import Typography from "@material-ui/core/Typography";

import { userId } from "./AppLayout/user";
import { MonsterGraphic } from "./Assets/Graphics";
import { LandingCard } from "./card";
import { useAppContext } from "./Contexts";

const StyledLanding = styled.div({
  backgroundColor: "#f9f9fe",
  height: "100%",
  "& .welcome": {
    display: "flex",
    backgroundColor: "white",
    padding: "32px 80px",
  },

  "& .welcome svg": {
    flex: "0 0 auto",
    marginRight: "24px",
  },

  "& .welcome .welcomeText": {
    flex: "1 1 auto",
  },

  "& .welcome .title": {
    fontWeight: "bold",
    fontSize: "22px",
    color: "#0d1030",
  },

  "& .welcome .subtitle": {
    fontSize: "16px",
    fontWeight: "normal",
    color: "rgba(13, 16, 48, 0.6)",
  },

  "& .content": {
    padding: "32px 80px",
  },
});

const Landing: React.FC<{}> = () => {
  const navigate = useNavigate();
  const { workflows } = useAppContext();
  const trendingWorkflows = [];
  workflows.forEach(workflow => {
    workflow.routes.forEach(route => {
      const title = route.displayName
        ? `${workflow.displayName}: ${route.displayName}`
        : workflow.displayName;
      if (route.trending) {
        trendingWorkflows.push({
          group: workflow.group,
          title,
          description: route.description,
          path: `${workflow.path}/${route.path}`,
        });
      }
    });
  });

  const navigateTo = (path: string) => {
    navigate(path);
  };

  return (
    <StyledLanding id="landing">
      <div className="welcome">
        <MonsterGraphic />
        <div className="welcomeText">
          <div className="title">Welcome {userId()}</div>
          <div className="subtitle">
            Clutch will assist you in safely modifying resources outside of the normal orchestration
            process.
          </div>
        </div>
      </div>
      <div className="content">
        {trendingWorkflows.length === 0 ? null : (
          <Grid container spacing={3} alignItems="center">
            <Grid item xs={12}>
              <Typography variant="h5">Trending Workflows</Typography>
            </Grid>
            {trendingWorkflows.map(workflow => (
              <Grid item xs={12} sm={12} md={6} lg={4} xl={4}>
                <LandingCard
                  group={workflow.group}
                  title={workflow.title}
                  description={workflow.description}
                  onClick={() => navigateTo(workflow.path)}
                  key={workflow.path}
                />
              </Grid>
            ))}
          </Grid>
        )}
      </div>
    </StyledLanding>
  );
};

export default Landing;
