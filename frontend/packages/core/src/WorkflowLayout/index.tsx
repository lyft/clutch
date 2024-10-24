import React from "react";
import { matchPath, Params, useParams } from "react-router";
import type { Interpolation } from "@emotion/styled";
import type { CSSObject, Theme } from "@mui/material";
import { alpha } from "@mui/system";

import type { Workflow } from "../AppProvider/workflow";
import Breadcrumbs from "../Breadcrumbs";
import { useLocation } from "../navigation";
import styled from "../styled";
import { Typography } from "../typography";
import { generateBreadcrumbsEntries } from "../utils";

export type LayoutVariant = "standard" | "wizard" | "custom";

// TODO: Define valid type variants
export type LayoutProps = {
  workflow: Workflow;
  variant?: LayoutVariant;
  title?: string | ((params: Params) => string);
  subtitle?: string;
  breadcrumbsOnly?: boolean;
  hideHeader?: boolean;
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
      padding: theme.spacing(theme.clutch.spacing.md),
    },
    wizard: {
      ...BASE_CONTAINER_STYLES,
      width: "800px", // Taken from the Wizard Component default width
      padding: theme.spacing(theme.clutch.spacing.lg, theme.clutch.spacing.none),
      margin: theme.spacing(theme.clutch.spacing.none, "auto"),
    },
    custom: {}, // No styles,
  };
  return layoutVariantStylesMap[variant];
};

const LayoutContainer = styled("div")(
  ({ $variant, theme }: StyledVariantComponentProps) =>
    getContainerVariantStyles($variant, theme) as Interpolation<void>
);

const PageHeader = styled("div")(({ $variant, theme }: StyledVariantComponentProps) => ({
  padding: theme.spacing(
    theme.clutch.spacing.none,
    $variant === "wizard" ? theme.clutch.spacing.md : theme.clutch.spacing.none
  ),
  paddingBottom: theme.spacing(theme.clutch.spacing.base),
  width: "100%",
}));

const PageHeaderBreadcrumbsWrapper = styled("div")(({ theme }: { theme: Theme }) => ({
  marginBottom: theme.spacing(theme.clutch.spacing.xs),
}));

const PageHeaderMainContainer = styled("div")(({ theme }: { theme: Theme }) => ({
  display: "flex",
  alignItems: "center",
  height: "70px",
  marginBottom: theme.spacing(theme.clutch.spacing.sm),
}));

const PageHeaderInformation = styled("div")({
  display: "flex",
  flexDirection: "column",
  justifyContent: "space-evenly",
  height: "100%",
});

const Title = styled(Typography)({
  lineHeight: 1,
});

const Subtitle = styled(Typography)(({ theme }: { theme: Theme }) => ({
  color: alpha(theme.colors.neutral[900], 0.45),
}));

const WorkflowLayout = ({
  workflow,
  variant = "standard",
  title = null,
  subtitle = null,
  breadcrumbsOnly = false,
  hideHeader = false,
  children,
}: React.PropsWithChildren<LayoutProps>) => {
  const params = useParams();
  const location = useLocation();
  const workflowPaths = workflow.routes.map(({ path }) => `/${workflow.path}/${path}`);
  const breadcrumbsEntries = generateBreadcrumbsEntries(
    location,
    (url: string) =>
      `/${workflow.path}` !== url &&
      !workflowPaths.includes(url) &&
      !workflowPaths.find(path => !!matchPath({ path }, url))
  );

  if (variant === "custom") {
    return <>{children}</>;
  }

  return (
    <LayoutContainer $variant={variant}>
      {!hideHeader && (
        <PageHeader $variant={variant}>
          <PageHeaderBreadcrumbsWrapper>
            <Breadcrumbs entries={breadcrumbsEntries} />
          </PageHeaderBreadcrumbsWrapper>
          {!breadcrumbsOnly && (title || subtitle) && (
            <PageHeaderMainContainer>
              <PageHeaderInformation>
                {title && (
                  <Title variant="h2" textTransform="capitalize">
                    {typeof title === "function" ? title(params) : title}
                  </Title>
                )}
                {subtitle && <Subtitle variant="subtitle2">{subtitle}</Subtitle>}
              </PageHeaderInformation>
            </PageHeaderMainContainer>
          )}
        </PageHeader>
      )}
      {children}
    </LayoutContainer>
  );
};

export default WorkflowLayout;
