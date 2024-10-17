import React from "react";
import { Interpolation } from "@emotion/styled";
import type { CSSObject, Theme } from "@mui/material";

import styled from "../styled";
import { Typography } from "../typography";

export type LayoutVariant = "standard" | "wizard" | "custom";

export type LayoutProps = {
  variant: LayoutVariant;
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
  variant,
  heading,
  hideHeader = false,
  children,
}: React.PropsWithChildren<LayoutProps>) => {
  return (
    <LayoutContainer $variant={variant}>
      {!hideHeader && (
        <PageHeader $variant={variant}>
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
