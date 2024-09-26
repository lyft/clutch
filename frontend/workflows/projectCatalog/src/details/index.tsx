import React from "react";
import { Grid, Tooltip, useTheme } from "@clutch-sh/core";
import CodeOffIcon from "@mui/icons-material/CodeOff";
import GroupIcon from "@mui/icons-material/Group";
import { useMediaQuery } from "@mui/material";
import { capitalize } from "lodash";

import type { CatalogDetailsChild, ProjectDetailsWorkflowProps } from "../types";

import { CardType, DynamicCard, MetaCard } from "./components/card";
import CatalogLayout from "./components/layout";
import { useProjectDetailsContext } from "./context";
import ProjectInfoCard from "./info";

const DisabledItem = ({ name }: { name: string }) => (
  <Grid item>
    <Tooltip title={`${capitalize(name)} is disabled`}>
      <CodeOffIcon />
    </Tooltip>
  </Grid>
);

const Details = ({ children, chips }: ProjectDetailsWorkflowProps) => {
  const { projectId, projectInfo } = useProjectDetailsContext() || {};
  const [metaCards, setMetaCards] = React.useState<CatalogDetailsChild[]>([]);
  const [dynamicCards, setDynamicCards] = React.useState<CatalogDetailsChild[]>([]);

  const theme = useTheme();
  const shrink = useMediaQuery(theme.breakpoints.down("lg"));

  React.useEffect(() => {
    if (children) {
      const tempMetaCards: CatalogDetailsChild[] = [];
      const tempDynamicCards: CatalogDetailsChild[] = [];

      React.Children.forEach(children, child => {
        if (React.isValidElement(child)) {
          const { type } = child?.props || {};

          switch (type) {
            case CardType.METADATA:
              tempMetaCards.push(child as CatalogDetailsChild);
              break;
            case CardType.DYNAMIC:
              tempDynamicCards.push(child as CatalogDetailsChild);
              break;
            default: {
              if (child.type === DynamicCard) {
                tempDynamicCards.push(child as CatalogDetailsChild);
              } else if (child.type === MetaCard) {
                tempMetaCards.push(child as CatalogDetailsChild);
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
    <>
      <Grid
        container
        item
        direction={shrink ? "row" : "column"}
        xs={12}
        sm={12}
        md={5}
        lg={4}
        xl={3}
        spacing={2}
        flexWrap="nowrap"
      >
        <Grid item>
          {projectInfo && (
            <MetaCard
              title={getOwner(projectInfo?.owners ?? []) || projectInfo?.name}
              titleIcon={<GroupIcon />}
              loadingIndicator={false}
              endAdornment={
                projectInfo?.data?.disabled ? (
                  <DisabledItem name={projectInfo?.name ?? projectId ?? ""} />
                ) : null
              }
            >
              {projectInfo && <ProjectInfoCard projectData={projectInfo} addtlChips={chips} />}
            </MetaCard>
          )}
        </Grid>
        {metaCards.length > 0 && metaCards.map(card => <Grid item>{card}</Grid>)}
      </Grid>
      <Grid
        container
        item
        direction="column"
        xs={12}
        sm={12}
        md={7}
        lg={8}
        xl={9}
        spacing={2}
        flexWrap="nowrap"
      >
        {dynamicCards.length > 0 && dynamicCards.map(card => <Grid item>{card}</Grid>)}
      </Grid>
    </>
  );
};

const CatalogDetailsPage = ({
  allowDisabled,
  configLinks,
  ...props
}: ProjectDetailsWorkflowProps) => {
  return (
    <CatalogLayout allowDisabled={allowDisabled} configLinks={configLinks ?? []}>
      <Details allowDisabled={allowDisabled} {...props} />
    </CatalogLayout>
  );
};

export default CatalogDetailsPage;
