import React from "react";
import { matchPath } from "react-router-dom";
import type { Interpolation } from "@emotion/styled";
import type { CSSObject, Theme } from "@mui/material";
import { alpha } from "@mui/material";

import type { Workflow } from "../AppProvider/workflow";
import Breadcrumbs from "../Breadcrumbs";
import Loadable from "../loading";
import { useLocation } from "../navigation";
import styled from "../styled";
import { Typography } from "../typography";
import { generateBreadcrumbsEntries } from "../utils";

import { useWorkflowLayoutContext } from "./context";

export type LayoutVariant = "standard" | "wizard";

export type LayoutProps = {
  workflow: Workflow;
  variant?: LayoutVariant | null;
  title?: string;
  subtitle?: string;
  breadcrumbsOnly?: boolean;
  hideHeader?: boolean;
  usesContext?: boolean;
};

type StyledVariantComponentProps = {
  theme: Theme;
  $variant: LayoutVariant;
};

const BASE_CONTAINER_STYLES: CSSObject = {
  display: "flex",
  flexDirection: "column",
  width: "100%",
  overflowY: "auto",
};

const getContainerVariantStyles = (variant: LayoutVariant, theme: Theme) => {
  const layoutVariantStylesMap: { [key in LayoutVariant]: CSSObject } = {
    standard: {
      ...BASE_CONTAINER_STYLES,
      padding: theme.spacing("md"),
    },
    wizard: {
      ...BASE_CONTAINER_STYLES,
      width: "800px", // Taken from the Wizard Component default width
      padding: theme.spacing("lg", "none"),
      margin: theme.spacing("none", "auto"),
    },
  };
  return layoutVariantStylesMap[variant];
};

const LayoutContainer = styled("div")(
  ({ $variant, theme }: StyledVariantComponentProps) =>
    getContainerVariantStyles($variant, theme) as Interpolation<void>
);

const PageHeader = styled("div")(({ $variant, theme }: StyledVariantComponentProps) => ({
  padding: theme.spacing("none", $variant === "wizard" ? "md" : "none"),
  paddingBottom: theme.spacing("base"),
  width: "100%",
}));

const PageHeaderBreadcrumbsWrapper = styled("div")(({ theme }: { theme: Theme }) => ({
  marginBottom: theme.spacing("xs"),
}));

const PageHeaderMainContainer = styled("div")({
  display: "flex",
  flexWrap: "wrap",
  justifyContent: "space-between",
  alignItems: "center",
  minHeight: "70px",
});

const PageHeaderInformation = styled("div")({
  display: "flex",
  flexDirection: "column",
  justifyContent: "space-evenly",
  height: "70px",
});

const PageHeaderSideContent = styled("div")({
  display: "flex",
  flexDirection: "column",
  justifyContent: "space-evenly",
  height: "70px",
});

const Title = styled(Typography)({
  lineHeight: 1,
  textTransform: "capitalize",
});

const Subtitle = styled(Typography)(({ theme }: { theme: Theme }) => ({
  color: alpha(theme.colors.neutral[900], 0.45),
  whiteSpace: "nowrap",
}));

const WorkflowLayout = ({
  workflow,
  variant = null,
  title = null,
  subtitle = null,
  breadcrumbsOnly = false,
  hideHeader = false,
  usesContext = false,
  children,
}: React.PropsWithChildren<LayoutProps>) => {
  const [headerLoading, setHeaderLoading] = React.useState(usesContext);

  const location = useLocation();
  const context = useWorkflowLayoutContext();

  const headerTitle = context?.title || title;
  const headerSubtitle = context?.subtitle || subtitle;

  React.useEffect(() => {
    if (context) {
      // Done to avoid a flash of the default title and subtitle
      setTimeout(() => setHeaderLoading(false), 750);
    }
  }, [context]);

  if (variant === null) {
    return <>{children}</>;
  }

  const workflowPaths = workflow.routes.map(({ path }) => `/${workflow.path}/${path}`);
  const breadcrumbsEntries = generateBreadcrumbsEntries(
    location,
    url => !!workflowPaths.find(path => !!matchPath({ path }, url))
  );

  return (
    <LayoutContainer $variant={variant}>
      {!hideHeader && (
        <PageHeader $variant={variant}>
          <PageHeaderBreadcrumbsWrapper>
            <Breadcrumbs entries={breadcrumbsEntries} />
          </PageHeaderBreadcrumbsWrapper>
          {!breadcrumbsOnly && (headerTitle || headerSubtitle) && (
            <PageHeaderMainContainer>
              <Loadable isLoading={headerLoading}>
                <PageHeaderInformation>
                  {headerTitle && <Title variant="h2">{headerTitle}</Title>}
                  {headerSubtitle && <Subtitle variant="subtitle2">{headerSubtitle}</Subtitle>}
                </PageHeaderInformation>
                {context?.headerContent && (
                  <PageHeaderSideContent>{context.headerContent}</PageHeaderSideContent>
                )}
              </Loadable>
            </PageHeaderMainContainer>
          )}
        </PageHeader>
      )}
      {children}
    </LayoutContainer>
  );
};

export default WorkflowLayout;
