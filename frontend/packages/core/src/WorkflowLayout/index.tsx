import React from "react";
import { useParams } from "react-router-dom";
import type { Interpolation } from "@emotion/styled";
import type { CSSObject, Theme } from "@mui/material";
import { alpha } from "@mui/material";

import type { Workflow } from "../AppProvider/workflow";
import Breadcrumbs from "../Breadcrumbs";
import { useLocation } from "../navigation";
import styled from "../styled";
import { Typography } from "../typography";
import { generateBreadcrumbsEntries } from "../utils";

export type LayoutVariant = "standard" | "wizard";

export type LayoutProps = {
  workflowsInPath: Array<Workflow>;
  variant?: LayoutVariant | null;
  title?: string | ((params: Record<string, string>) => string);
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

const PageHeaderMainContainer = styled("div")(({ theme }: { theme: Theme }) => ({
  display: "flex",
  alignItems: "center",
  height: "70px",
  marginBottom: theme.spacing("sm"),
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
  workflowsInPath,
  variant = null,
  title = null,
  subtitle = null,
  breadcrumbsOnly = false,
  hideHeader = false,
  children,
}: React.PropsWithChildren<LayoutProps>) => {
  const params = useParams();
  const location = useLocation();

  const entries = generateBreadcrumbsEntries(workflowsInPath, location);

  if (variant === null) {
    return <>{children}</>;
  }

  return (
    <LayoutContainer $variant={variant}>
      {!hideHeader && (
        <PageHeader $variant={variant}>
          <PageHeaderBreadcrumbsWrapper>
            <Breadcrumbs entries={entries} />
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
