import { Location, matchPath } from "react-router-dom";

import type { Workflow } from "../AppProvider/workflow";
import type { BreadcrumbEntry } from "../Breadcrumbs";

const HOME_ENTRY = { label: "Home", url: "/" };

const generateBreadcrumbsEntries = (workflowsInPath: Array<Workflow>, location: Location) => {
  // The first workflow in the will contain
  // the same path and displayName as the others
  const firstWorkflow = workflowsInPath[0];

  if (!firstWorkflow) {
    return [HOME_ENTRY];
  }

  // Get a single level list of the routes available
  const allRoutes = workflowsInPath.flatMap(w => w.routes);

  // Add to every item in the routes list the workflow path prefix
  const fullPaths = allRoutes.map(({ path }) => `/${firstWorkflow.path}/${path}`);

  // Generate a list of path segments from the location
  const pathSegments = decodeURIComponent(location.pathname)
    .split("/")
    .slice(1, location.pathname.endsWith("/") ? -1 : undefined); // in case of a trailing `/`

  const entries: Array<BreadcrumbEntry> = [HOME_ENTRY].concat(
    pathSegments.map((segment, index) => {
      const nextIndex = index + 1;
      const url = `/${pathSegments.slice(0, nextIndex).join("/")}`;

      const path = fullPaths.find(p => !!matchPath(p, url));

      // If there is a matched path, it's used to find the route that contains its displayName
      const route = path
        ? allRoutes.find(r =>
            r.path.startsWith("/")
              ? r.path
              : // Done in case of an empty path or missing a leading `/`
                `/${r.path}` === `/${path.split("/").slice(2).join("/")}`
          )
        : null;

      return {
        // For the label:
        // - Prioritize the display name
        // - Handle the case of a single route with an unusual long name
        // - Default to the path segment
        label:
          route?.displayName || (allRoutes.length === 1 && firstWorkflow.displayName) || segment,
        // Set a null url if there is no path or for the last segment
        url: !!path && pathSegments.length !== nextIndex ? url : null,
      };
    })
  );

  return entries;
};

export default generateBreadcrumbsEntries;
