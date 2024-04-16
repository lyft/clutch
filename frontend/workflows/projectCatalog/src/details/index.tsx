import React from "react";
import { useParams } from "react-router-dom";
import type { clutch as IClutch } from "@clutch-sh/api";
import {
  checkFeatureEnabled,
  Grid,
  IconButton,
  Link,
  Popper,
  PopperItem,
  QuickLinkGroup,
  QuickLinksCard,
  styled,
  Tooltip,
  Typography,
} from "@clutch-sh/core";
import GroupIcon from "@mui/icons-material/Group";
import SettingsIcon from "@mui/icons-material/Settings";
import { capitalize, isEmpty } from "lodash";
import CodeOffIcon from "@mui/icons-material/CodeOff";

import type { CatalogDetailsChild, ProjectConfigLink, ProjectDetailsWorkflowProps } from "..";

import { CardType, DynamicCard, MetaCard } from "./card";
import { ProjectDetailsContext } from "./context";
import ProjectHeader from "./header";
import { fetchProjectInfo } from "./helpers";
import ProjectInfoCard from "./info";

interface QuickLinksAndSettingsProps {
  linkGroups: QuickLinkGroup[];
  configLinks?: ProjectConfigLink[];
}

const StyledContainer = styled(Grid)({
  padding: "16px",
});

const StyledHeadingContainer = styled(Grid)({
  marginBottom: "24px",
});

const StyledPopperItem = styled(PopperItem)({
  "&&&": {
    height: "auto",
  },
  "& span.MuiTypography-root": {
    padding: "0",
  },
  "& a.MuiTypography-root": {
    padding: "4px 16px",
  },
});

const DisabledItem = ({ name }: { name: string }) => (
  <Grid item>
    <Tooltip title={`${capitalize(name)} is disabled`}>
      <CodeOffIcon />
    </Tooltip>
  </Grid>
);

const QuickLinksAndSettingsBtn = ({ linkGroups, configLinks = [] }: QuickLinksAndSettingsProps) => {
  const { projectId } = useParams();
  const anchorRef = React.useRef(null);
  const [open, setOpen] = React.useState(false);
  const [links, setLinks] = React.useState<ProjectConfigLink[]>(configLinks);

  React.useEffect(() => {
    const projectConfigFlag = checkFeatureEnabled({ feature: "projectCatalogSettings" });
    if (projectConfigFlag) {
      setLinks([
        {
          title: "Project Configuration",
          path: `/catalog/${projectId}/config`,
          icon: <SettingsIcon fontSize="small" />,
        },
        ...links,
      ]);
    }
  }, []);

  return (
    <Grid
      container
      direction="row"
      alignItems="center"
      justifyContent="flex-end"
      spacing={1}
      style={{
        padding: "8px 0px 0px 0px",
      }}
    >
      {!isEmpty(linkGroups) && (
        <Grid item>
          <QuickLinksCard linkGroups={linkGroups} />
        </Grid>
      )}
      {links && links.length > 0 && (
        <Grid item>
          <IconButton ref={anchorRef} onClick={() => setOpen(o => !o)} size="medium">
            <SettingsIcon />
          </IconButton>
          <Popper
            open={open}
            anchorRef={anchorRef}
            onClickAway={() => setOpen(false)}
            placement="bottom-end"
          >
            {links.map(link => (
              <StyledPopperItem key={link.title}>
                <Link href={link.path}>
                  <Grid container gap={0.5}>
                    {link.icon && <Grid item>{link.icon}</Grid>}
                    <Grid item>
                      <Typography variant="body2" color="inherit">
                        {link.title}
                      </Typography>
                    </Grid>
                  </Grid>
                </Link>
              </StyledPopperItem>
            ))}
          </Popper>
        </Grid>
      )}
    </Grid>
  );
};

const Details: React.FC<ProjectDetailsWorkflowProps> = ({
  children,
  chips,
  allowDisabled,
  configLinks = [],
}) => {
  const { projectId } = useParams();
  const [projectInfo, setProjectInfo] = React.useState<IClutch.core.project.v1.IProject | null>(
    null
  );
  const [metaCards, setMetaCards] = React.useState<CatalogDetailsChild[]>([]);
  const [dynamicCards, setDynamicCards] = React.useState<CatalogDetailsChild[]>([]);
  const projInfo = React.useMemo(() => ({ projectInfo }), [projectInfo]);

  React.useEffect(() => {
    if (children) {
      const tempMetaCards: CatalogDetailsChild[] = [];
      const tempDynamicCards: CatalogDetailsChild[] = [];

      React.Children.forEach(children, child => {
        if (React.isValidElement(child)) {
          const { type } = child?.props || {};

          switch (type) {
            case CardType.METADATA:
              tempMetaCards.push(child);
              break;
            case CardType.DYNAMIC:
              tempDynamicCards.push(child);
              break;
            default: {
              if (child.type === DynamicCard) {
                tempDynamicCards.push(child);
              } else if (child.type === MetaCard) {
                tempMetaCards.push(child);
              }
              // Do nothing, invalid card
            }
          }
        }
      });

      setMetaCards(tempMetaCards);
      setDynamicCards(tempDynamicCards);
    }
  }, [children]);

  /**
   * Takes an array of owner emails and returns the first one capitalized
   * (ex: ["clutch-team@lyft.com"] -> "Clutch Team")
   */
  const getOwner = (owners: string[]): string => {
    if (owners && owners.length) {
      const firstOwner = owners[0];

      return firstOwner
        .replace(/@.*\..*/, "")
        .split("-")
        .join(" ");
    }

    return "";
  };

  return (
    <ProjectDetailsContext.Provider value={projInfo}>
      <StyledContainer container direction="row" wrap="nowrap">
        {/* Column for project details and header */}
        <Grid container item direction="column" xs={12} sm={12} md={12} lg={12} xl={12}>
          <Grid container item>
            <StyledHeadingContainer item xs={6} sm={6} md={7} lg={8} xl={9}>
              {/* Static Header */}
              <ProjectHeader
                title={projectId}
                routes={[{ title: projectId }]}
                description={projectInfo?.data?.description as string}
              />
            </StyledHeadingContainer>
            {projectInfo && (
              <Grid container item xs={12} sm={12} md={5} lg={4} xl={3} spacing={2}>
                <QuickLinksAndSettingsBtn
                  linkGroups={(projectInfo.linkGroups as QuickLinkGroup[]) || []}
                  configLinks={configLinks ?? []}
                />
              </Grid>
            )}
          </Grid>
          <Grid container item direction="row" spacing={2}>
            <Grid container item xs={12} sm={12} md={5} lg={4} xl={3} spacing={2}>
              <Grid item xs={12}>
                {/* Static Info Card */}
                <MetaCard
                  title={getOwner(projectInfo?.owners ?? []) || projectId}
                  titleIcon={<GroupIcon />}
                  fetchDataFn={() => fetchProjectInfo(projectId, allowDisabled)}
                  onSuccess={(data: unknown) =>
                    setProjectInfo(data as IClutch.core.project.v1.IProject)
                  }
                  autoReload
                  loadingIndicator={false}
                  endAdornment={
                    projectInfo?.data?.disabled ? <DisabledItem name={projectId} /> : null
                  }
                >
                  {projectInfo && <ProjectInfoCard projectData={projectInfo} addtlChips={chips} />}
                </MetaCard>
              </Grid>
              {/* Custom Meta Cards */}
              {metaCards.length > 0 &&
                metaCards.map(card => (
                  <Grid item xs={12}>
                    {card}
                  </Grid>
                ))}
            </Grid>
            <Grid container item xs={12} sm={12} md={7} lg={8} xl={9} spacing={1}>
              {/* Custom Dynamic Cards */}
              {dynamicCards.length > 0 &&
                dynamicCards.map(card => (
                  <Grid item xs={12} key={`dynamic-${card.key}`}>
                    {card}
                  </Grid>
                ))}
            </Grid>
          </Grid>
        </Grid>
      </StyledContainer>
    </ProjectDetailsContext.Provider>
  );
};

export default Details;
