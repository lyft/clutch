import type { Location } from "history";

import type { BreadcrumbEntry } from "../Breadcrumbs";

const generateBreadcrumbsEntries = (location: Location, validateUrl: (url: string) => boolean) => {
  const labels = decodeURIComponent(location.pathname)
    .split("/")
    .slice(1, location.pathname.endsWith("/") ? -1 : undefined);

  const entries: Array<BreadcrumbEntry> = [{ label: "Home", url: "/" }].concat(
    labels.map((label, index) => {
      let url = `/${labels.slice(0, index + 1).join("/")}`;

      if (validateUrl(url)) {
        url = undefined;
      }

      return {
        label,
        url,
      };
    })
  );

  return entries;
};

export default generateBreadcrumbsEntries;
