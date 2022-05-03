import React from "react";
import { useLocation, useNavigate, useParams } from "react-router-dom";
import type { clutch as IClutch } from "@clutch-sh/api";
import type { ClutchError } from "@clutch-sh/core";
import { Error, Grid, styled, Tab, Tabs } from "@clutch-sh/core";

import type { ProjectConfigPage, ProjectDetailsConfigWorkflowProps } from "..";
import { ProjectDetailsContext } from "../details/context";
import ProjectHeader from "../details/header";
import { fetchProjectInfo } from "../details/helpers";

const StyledContainer = styled(Grid)({
  padding: "32px",
});

const Config: React.FC<ProjectDetailsConfigWorkflowProps> = ({ children, defaultRoute }) => {
  const { projectId, configType = defaultRoute } = useParams();
  const location = useLocation();
  const navigate = useNavigate();
  const [projectInfo, setProjectInfo] = React.useState<IClutch.core.project.v1.IProject | null>(
    null
  );
  const [error, setError] = React.useState<ClutchError | null>(null);
  const [configPages, setConfigPages] = React.useState<ProjectConfigPage[]>([]);
  const [selectedPage, setSelectedPage] = React.useState<number>(0);

  React.useEffect(() => {
    fetchProjectInfo(projectId).then(setProjectInfo).catch(setError);
  }, []);

  React.useEffect(() => {
    if (configPages && configPages.length) {
      const splitLoc = location.pathname.split("/");
      if (splitLoc[splitLoc.length - 1] !== "config") {
        splitLoc.splice(splitLoc.length - 1, 1, configPages[selectedPage]?.props?.path);
      } else {
        splitLoc.push(configPages[selectedPage]?.props?.path);
      }
      navigate(splitLoc.join("/"));
    }
  }, [configPages, selectedPage]);

  React.useEffect(() => {
    if (children) {
      const validPages: ProjectConfigPage[] = [];

      React.Children.forEach(children, (child, index) => {
        if (React.isValidElement(child)) {
          const { title, path } = child?.props;

          if (title) {
            validPages.push(child);

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
    <ProjectDetailsContext.Provider value={{ projectInfo }}>
      <StyledContainer container spacing={2}>
        <Grid item xs={12}>
          <ProjectHeader
            title="Project Configuration"
            routes={[
              { title: "Details", path: `${projectId}` },
              { title: "Project Configuration" },
              { title: configType || defaultRoute },
            ]}
            description="Edit your projects' settings."
          />
        </Grid>
        {configPages && configPages.length > 1 ? (
          <Grid item xs={12}>
            <Tabs value={selectedPage} centered>
              {configPages.map((page, i) => (
                <Tab label={page.props.title} onClick={() => setSelectedPage(i)} />
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
          {configPages &&
            configPages.length > 0 &&
            React.cloneElement(configPages[selectedPage], {
              ...configPages[selectedPage].props,
              onError: setError,
            })}
        </Grid>
      </StyledContainer>
    </ProjectDetailsContext.Provider>
  );
};

export default Config;
