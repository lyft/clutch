import React from "react";
import { matchPath } from "react-router";
import type { Interpolation } from "@emotion/styled";
import type { CSSObject, Theme } from "@mui/material";

import { Workflow } from "../AppProvider/workflow";
import Breadcrumbs, { BreadcrumbEntry } from "../Breadcrumbs";
import { useLocation } from "../navigation";
import styled from "../styled";
import { Typography } from "../typography";

export type LayoutVariant = "standard" | "wizard" | "custom";

export type LayoutProps = {
  workflow: Workflow;
  variant?: LayoutVariant;
  heading?: string | React.ReactElement;
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
    // No styles
    custom: {},
  };
  return layoutVariantStylesMap[variant];
};

const LayoutContainer = styled("div")(
  ({ $variant, theme }: StyledVariantComponentProps) =>
    getContainerVariantStyles($variant, theme) as Interpolation<void>
);

const PageHeader = styled("div")(({ $variant, theme }: StyledVariantComponentProps) => ({
  padding: theme.spacing(
    theme.clutch.spacing.base,
    $variant === "wizard" ? theme.clutch.spacing.md : theme.clutch.spacing.none
  ),
  width: "100%",
}));

const HeaderTitle = styled(Typography)({
  lineHeight: 1,
});

const WorkflowLayout = ({
  workflow,
  variant = "standard",
  heading = null,
  hideHeader = false,
  children,
}: React.PropsWithChildren<LayoutProps>) => {
  const location = useLocation();
  const workflowPaths = workflow.routes.map(({ path }) => `/${workflow.path}/${path}`);

  const labels = location.pathname
    .split("/")
    .slice(1, location.pathname.endsWith("/") ? -1 : undefined);

  const breadcrumbsEntries: Array<BreadcrumbEntry> = [{ label: "Home", url: "/" }].concat(
    labels.map((label, index) => {
      let url = `/${labels.slice(0, index + 1).join("/")}`;
      const isInvalidUrl =
        `/${workflow.path}` !== url &&
        !workflowPaths.includes(url) &&
        !workflowPaths.find(path => !!matchPath({ path }, url));

      if (isInvalidUrl) {
        url = undefined;
      }

      return {
        label,
        url,
      };
    })
  );

  return (
    <LayoutContainer $variant={variant}>
      {!hideHeader && (
        <PageHeader $variant={variant}>
          <Breadcrumbs entries={breadcrumbsEntries} />
          {heading && (
            <>
              {React.isValidElement(heading) ? (
                heading
              ) : (
                <HeaderTitle variant="h2">{heading}</HeaderTitle>
              )}
            </>
          )}
        </PageHeader>
      )}
      {children}
    </LayoutContainer>
  );
};

export default WorkflowLayout;
