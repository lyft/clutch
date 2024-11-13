import { Location, matchPath } from "react-router-dom";

import type { Workflow } from "../AppProvider/workflow";
import type { BreadcrumbEntry } from "../Breadcrumbs";

const HOME_ENTRY = { label: "Home", url: "/" };

const generateBreadcrumbsEntries = (workflowsInPath: Array<Workflow>, location: Location) => {
  const firstWorkflow = workflowsInPath[0];
  const allRoutes = workflowsInPath.flatMap(w => w.routes);
  const fullPaths = allRoutes.map(({ path }) => `/${firstWorkflow.path}/${path}`);

  const labels = decodeURIComponent(location.pathname)
    .split("/")
    .slice(1, location.pathname.endsWith("/") ? -1 : undefined);

  const entries: Array<BreadcrumbEntry> = [HOME_ENTRY].concat(
    labels.map((defaultLabel, index) => {
      const nextIndex = index + 1;
      const url = `/${labels.slice(0, nextIndex).join("/")}`;

      const path = fullPaths.find(p => !!matchPath(p, url));
      const route = path
        ? allRoutes.find(r =>
            r.path.startsWith("/")
              ? r.path
              : `/${r.path}` === `/${path.split("/").slice(2).join("/")}`
          )
        : null;

      return {
        label:
          route?.displayName ||
          (allRoutes.length === 1 && firstWorkflow.displayName) ||
          defaultLabel,
        url: !!path && labels.length !== nextIndex ? url : null,
      };
    })
  );

  return entries;
};

export default generateBreadcrumbsEntries;
