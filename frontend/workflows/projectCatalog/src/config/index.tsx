import React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";
import type { ClutchError } from "@clutch-sh/core";
import {
  Error,
  Grid,
  styled,
  Tab,
  Tabs,
  useLocation,
  useNavigate,
  useParams,
} from "@clutch-sh/core";

import type { ProjectConfigPage, ProjectDetailsConfigWorkflowProps } from "..";
import { ProjectDetailsContext } from "../details/context";
import ProjectHeader from "../details/header";
import { fetchProjectInfo } from "../details/helpers";

const StyledContainer = styled(Grid)({
  padding: "32px",
});

const Config: React.FC<ProjectDetailsConfigWorkflowProps> = ({ children, defaultRoute = "/" }) => {
  const { projectId, configType = defaultRoute } = useParams();
  const location = useLocation();
  const navigate = useNavigate();
  const [projectInfo, setProjectInfo] = React.useState<IClutch.core.project.v1.IProject | null>(
    null
  );
  const [error, setError] = React.useState<ClutchError | null>(null);
  const [configPages, setConfigPages] = React.useState<ProjectConfigPage[]>([]);
  const [selectedPage, setSelectedPage] = React.useState<number>(0);
  const projInfo = React.useMemo(() => ({ projectInfo }), [projectInfo]);

  React.useEffect(() => {
    fetchProjectInfo(projectId).then(setProjectInfo).catch(setError);
  }, []);

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
    <ProjectDetailsContext.Provider value={projInfo}>
      <StyledContainer container spacing={2}>
        <Grid item xs={12}>
          <ProjectHeader
            title={`${projectInfo?.name ?? projectId} Configuration`}
            routes={[
              { title: projectId, path: `${projectId}` },
              { title: "Configuration" },
              { title: configType || defaultRoute },
            ]}
            description="Edit your projects' settings. Some changes may require additional approval from another owner of the project."
          />
        </Grid>
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
      </StyledContainer>
    </ProjectDetailsContext.Provider>
  );
};

export default Config;
