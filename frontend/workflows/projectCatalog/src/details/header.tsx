import React from "react";
import { Grid, styled, Typography, useTheme } from "@clutch-sh/core";
import { alpha, Breadcrumbs, Link } from "@mui/material";

interface Route {
  title: string;
  path?: string;
}

interface BreadCrumbProps {
  routes: Route[];
}

interface ProjectHeaderProps extends BreadCrumbProps {
  title: string;
  description?: string;
}

const StyledHeading = styled("div")({
  padding: "8px 0px 8px 0px",
  textTransform: "capitalize",
});

const StyledContainer = styled(Grid)({
  width: "100%",
  height: "100%",
});

const StyledCrumb = styled(Typography)({
  textTransform: "uppercase",
});

const BreadCrumbs = ({ routes = [] }: BreadCrumbProps) => {
  const theme = useTheme();
  routes.unshift({ title: "Project Catalog", path: "/catalog" });

  let builtRoute = routes[0].path;

  const buildCrumb = (route: Route) => {
    if (route?.path && route?.path !== builtRoute) {
      builtRoute += `/${route.path}`;
    }

    return (
      <StyledCrumb
        variant="caption2"
        color={alpha(theme.palette.secondary[900], 0.48)}
        key={route.title}
      >
        {route.path ? (
          <Link color="inherit" href={builtRoute} underline="hover">
            {route.title}
          </Link>
        ) : (
          route.title
        )}
      </StyledCrumb>
    );
  };

  return <Breadcrumbs aria-label="breadcrumbs">{routes.map(buildCrumb)}</Breadcrumbs>;
};

const ProjectHeader = ({ title, routes, description = "" }: ProjectHeaderProps) => (
  <StyledContainer container direction="column">
    <Grid container item direction="row" alignItems="flex-end">
      <BreadCrumbs routes={routes} />
    </Grid>
    <StyledHeading>
      <Typography variant="h2">{title}</Typography>
    </StyledHeading>
    {description.length > 0 && <Typography variant="body2">{description}</Typography>}
  </StyledContainer>
);

export default ProjectHeader;
