import React from "react";
import { Typography, useTheme } from "@clutch-sh/core";
import { alpha, Breadcrumbs, Link } from "@mui/material";

interface Route {
  title: string;
  path?: string;
}

export interface BreadCrumbsProps {
  routes?: Route[];
}

const BreadCrumbs = ({ routes = [] }: BreadCrumbsProps) => {
  const theme = useTheme();
  routes.unshift({ title: "Project Catalog", path: "/catalog" });

  let builtRoute = routes[0].path;

  const buildCrumb = (route: Route) => {
    if (route?.path && route?.path !== builtRoute) {
      builtRoute += `/${route.path}`;
    }

    return (
      <Typography
        textTransform="uppercase"
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
      </Typography>
    );
  };

  return <Breadcrumbs aria-label="breadcrumbs">{routes.map(buildCrumb)}</Breadcrumbs>;
};

export default BreadCrumbs;
