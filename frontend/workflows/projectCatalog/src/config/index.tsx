import React from "react";
import { useParams } from "react-router-dom";
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
            description="Edit your projects' settings. Changes will immediately take effect on save. Some changes may require additional approval from another owner of the project"
          />
        </Grid>
        {configPages && configPages.length > 0 ? (
          <Grid item xs={12}>
            <Tabs>
              {configPages.map((page, i) => (
                <Tab label={page.props.title} onClick={_ => setSelectedPage(i)} />
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
