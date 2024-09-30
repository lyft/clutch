import React from "react";
import type { ClutchError } from "@clutch-sh/core";
import { Error, Grid, Tab, Tabs, useLocation, useNavigate, useParams } from "@clutch-sh/core";

import type { ProjectCatalogProps, WorkflowProps } from "../../types";
import CatalogLayout from "../components/layout";

export interface ProjectConfigProps {
  title: string;
  path: string;
  onError?: (error: ClutchError) => void;
}

export type ProjectConfigPage = React.ReactElement<ProjectConfigProps>;

export interface ProjectDetailsConfigWorkflowProps extends WorkflowProps, ProjectCatalogProps {
  children?: ProjectConfigPage | ProjectConfigPage[];
  // eslint-disable-next-line react/no-unused-prop-types
  description?: string;
  defaultRoute?: string;
}

const Config = ({ children, defaultRoute = "/" }: ProjectDetailsConfigWorkflowProps) => {
  const { configType = defaultRoute } = useParams();
  const location = useLocation();
  const navigate = useNavigate();
  const [error, setError] = React.useState<ClutchError | null>(null);
  const [configPages, setConfigPages] = React.useState<ProjectConfigPage[]>([]);
  const [selectedPage, setSelectedPage] = React.useState<number>(0);

  React.useEffect(() => {
    if (configPages && configPages.length) {
      const splitLoc = location.pathname.split("/");
      const selectedPath = configPages[selectedPage]?.props?.path;

      if (selectedPath) {
        if (splitLoc[splitLoc.length - 1] !== "config") {
          splitLoc.splice(splitLoc.length - 1, 1, selectedPath);
        } else {
          splitLoc.push(selectedPath);
        }

        // Used to reduce the number of navigation calls when the user is navigating between tabs
        if (splitLoc.join("/") !== location.pathname.replace(/%20/, " ")) {
          navigate(
            {
              pathname: splitLoc.join("/"),
              search: window.location.search,
            },
            { replace: true }
          );
        }
      }
    }
  }, [configPages, selectedPage]);

  React.useEffect(() => {
    if (children) {
      const validPages: ProjectConfigPage[] = [];

      React.Children.forEach(children, (child, index) => {
        if (React.isValidElement(child)) {
          const { title, path, onError } = child?.props || {};

          if (title) {
            validPages.push(
              React.cloneElement(child, {
                onError: (e: ClutchError) => {
                  if (onError) {
                    onError(e);
                  }
                  setError(e);
                },
              })
            );

            if (configType === path) {
              setSelectedPage(index);
            }
          }
        }
      });

      setConfigPages(validPages);
    }
  }, [children]);

  return (
    <>
      {configPages && configPages.length > 1 ? (
        <Grid item xs={12}>
          <Tabs value={selectedPage} centered>
            {configPages.map((page, i) => (
              <Tab
                key={page.props.title}
                label={page.props.title}
                onClick={() => setSelectedPage(i)}
              />
            ))}
          </Tabs>
        </Grid>
      ) : null}
      {error && (
        <Grid item xs={12}>
          <Error subject={error} />
        </Grid>
      )}
      <Grid item xs={12}>
        {configPages && configPages.length > 0 && configPages[selectedPage]}
      </Grid>
    </>
  );
};

const CatalogConfigPage = ({
  defaultRoute = "/",
  description = "",
  ...props
}: ProjectDetailsConfigWorkflowProps) => {
  const { configType = defaultRoute } = useParams();

  return (
    <CatalogLayout
      title="Configuration"
      description={description}
      routes={[{ title: "Configuration" }, { title: configType || defaultRoute }]}
      allowDisabled={false}
      quickLinkSettings={false}
      configLinks={[]}
    >
      <Config defaultRoute={defaultRoute} {...props} />
    </CatalogLayout>
  );
};

export default CatalogConfigPage;
