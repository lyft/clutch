import React from "react";
import styled from "@emotion/styled";
import { alpha, Grid, Theme } from "@mui/material";
import Typography from "@mui/material/Typography";

import { userId } from "./AppLayout/user";
import { workflowsByTrending } from "./AppLayout/utils";
import { MonsterGraphic } from "./Assets/Graphics";
import { LandingCard } from "./card";
import { useAppContext } from "./Contexts";
import { useNavigate } from "./navigation";

const StyledLanding = styled.div(({ theme }: { theme: Theme }) => ({
  display: "flex",
  flexDirection: "column",
  flexGrow: 1,
  "& .welcome": {
    display: "flex",
    backgroundColor:
      theme.palette.mode === "light" ? theme.palette.contrastColor : theme.palette.background.paper,
    padding: "32px 80px",
  },

  "& .welcome svg": {
    flex: "0 0 auto",
    margin: "auto 24px auto 0",
  },

  "& .welcome .welcomeText": {
    flex: "1 1 auto",
  },

  "& .welcome .title": {
    fontWeight: "bold",
    fontSize: "22px",
    color: theme.palette.text.primary[900],
  },

  "& .welcome .subtitle": {
    fontSize: "16px",
    fontWeight: "normal",
    color: alpha(theme.palette.secondary[900], 0.6),
  },

  "& .content": {
    padding: "32px 80px",
  },
}));

const Landing: React.FC<{}> = () => {
  const navigate = useNavigate();
  const { workflows } = useAppContext();
  const trendingWorkflows = workflowsByTrending(workflows);

  const navigateTo = (path: string) => {
    navigate(path);
  };

  return (
    <StyledLanding id="landing">
      <div className="welcome">
        <MonsterGraphic />
        <div className="welcomeText">
          <div className="title">Welcome&nbsp;{userId()}</div>
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
              <Grid key={workflow.path} item xs={12} sm={12} md={6} lg={4} xl={4}>
                <LandingCard
                  group={workflow.group}
                  title={workflow.displayName}
                  description={workflow.description}
                  onClick={() => navigateTo(workflow.path)}
                  key={workflow.path}
                  icon={workflow.icon}
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
