import React from "react";
import type { CSSObject, Theme } from "@mui/material";

import styled from "../styled";
import { Typography } from "../typography";

type LayoutVariant = "standard" | "wizard";

type LayoutProps = {
  variant: LayoutVariant;
  heading?: string | React.ReactElement;
  showHeader?: boolean;
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
  };
  return layoutVariantStylesMap[variant];
};

const LayoutContainer = styled("div")(
  ({ $variant, theme }: StyledVariantComponentProps) =>
    getContainerVariantStyles($variant, theme) as any
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
  showHeader,
  heading,
  children,
}: React.PropsWithChildren<LayoutProps>) => {
  return (
    <LayoutContainer $variant={variant}>
      {showHeader !== false && (
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
