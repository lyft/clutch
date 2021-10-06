import * as React from "react";
import styled from "@emotion/styled";
import {
  Avatar,
  Card as MuiCard,
  CardActionArea,
  CardActionAreaProps,
  Divider,
  Grid,
} from "@material-ui/core";
import KeyboardArrowDownIcon from "@material-ui/icons/KeyboardArrowDown";
import KeyboardArrowUpIcon from "@material-ui/icons/KeyboardArrowUp";
import type { SpacingProps as MuiSpacingProps } from "@material-ui/system";
import { spacing } from "@material-ui/system";

import { IconButton } from "./button";
import { Typography, TypographyProps } from "./typography";

// TODO: seperate out the different card parts into various files

const StyledCard = styled(MuiCard)({
  boxShadow: "0px 4px 6px rgba(53, 72, 212, 0.2)",
  border: "1px solid rgba(13, 16, 48, 0.1)",
});

export interface CardProps {
  children?: React.ReactNode | React.ReactNode[];
}

const Card = ({ children, ...props }: CardProps) => <StyledCard {...props}>{children}</StyledCard>;

const StyledCardHeaderContainer = styled.div({
  background: "#EBEDFB",
});

const StyledCardHeader = styled(Grid)({
  padding: "6px 8px",
  minHeight: "48px",
  ".MuiGrid-item": {
    padding: "0px 8px",
  },
});

const StyledCardHeaderAvatarContainer = styled.div({
  padding: "8px",
  height: "32px",
  width: "32px",
  alignSelf: "center",
  display: "flex",
});

// TODO: use material ui avatar component and implement figma design
const StyledCardHeaderAvatar = styled.div({
  width: "24px",
  height: "24px",
  fontSize: "18px",
  alignSelf: "center",
});

// TODO: make the divider a core component
const StyledDivider = styled(Divider)({
  color: "#A3A4B0",
  height: "24px",
  alignSelf: "center",
});

const StyledGridItem = styled(Grid)({
  textAlign: "center",
});

export interface CardHeaderSummaryProps {
  title: React.ReactElement<TypographyProps> | React.ReactElement<TypographyProps>[];
  subheader?: string;
}

interface CardHeaderProps {
  actions?: React.ReactNode;
  avatar: React.ReactNode;
  children?: React.ReactNode;
  summary?: CardHeaderSummaryProps[];
  title: React.ReactNode;
}

const CardHeader = ({ actions, avatar, children, title, summary = [] }: CardHeaderProps) => {
  return (
    <StyledCardHeaderContainer>
      <StyledCardHeader container wrap="nowrap" alignItems="center">
        <StyledCardHeaderAvatarContainer>
          <StyledCardHeaderAvatar>{avatar}</StyledCardHeaderAvatar>
        </StyledCardHeaderAvatarContainer>
        <Grid container wrap="nowrap" alignItems="center">
          <Grid item xs>
            <Typography variant="h4">{title}</Typography>
          </Grid>
          {summary.map(section => (
            <>
              <StyledDivider orientation="vertical" />
              <StyledGridItem item xs>
                {section.title}
                {section.subheader && (
                  <Typography variant="body4" color="rgba(13, 16, 48, 0.6)">
                    {section.subheader}
                  </Typography>
                )}
              </StyledGridItem>
            </>
          ))}
        </Grid>
        {actions}
      </StyledCardHeader>
      {children}
    </StyledCardHeaderContainer>
  );
};

// Material UI Spacing system supports many props https://material-ui.com/system/spacing/#api
// We can add more to this list as use cases arise
interface SpacingProps extends Pick<MuiSpacingProps, "padding" | "p"> {}

const BaseCardContent = styled.div<SpacingProps>`
  ${spacing}
`;

const StyledCardContentContainer = styled.div((props: { maxHeight: number | "none" }) => ({
  "> .MuiPaper-root": {
    border: "0",
    borderRadius: "0",
  },
  overflow: "hidden",
  maxHeight: props.maxHeight,
}));

const BaseCardActionArea = styled(CardActionArea)<SpacingProps>`
  ${spacing}
`;

const StyledCardActionArea = styled(BaseCardActionArea)({
  ":hover": {
    backgroundColor: "#F5F6FD",
  },

  ":active": {
    backgroundColor: "#D7DAF6",
  },
});

const StyledExpandButton = styled(IconButton)({
  width: "32px",
  height: "32px",
  color: "#3548D4",
  ":hover": {
    backgroundColor: "transparent",
  },
});

interface CollapsibleState {
  /** The text to show in the collapse action area when the card content is collapsed/not collapsed.
   * By default, this is "See Less" for true and "See More" for false.
   */
  title: string;
  /** The icon to show in the collapse action area when the card content is collapsed/not collapsed.
   * By default, this is an ArrowUp icon for true and an ArrowDown icon for false.
   */
  icon: React.ReactElement;
}

interface CardContentCollapsibleProps {
  open?: CollapsibleState;
  closed?: CollapsibleState;
}

interface CardContentProps extends SpacingProps {
  children?: React.ReactNode | React.ReactNode[];
  /** Make the card content collapse at a set max height. The default is false. */
  collapsible?: boolean;
  /** The props for the collapse action */
  collapseAction?: CardContentCollapsibleProps;
  /** The max height of the card content. The default is none. */
  maxHeight?: number | "none";
}

const CardContent = ({
  children,
  collapsible = false,
  collapseAction = {
    open: { title: "See Less", icon: <KeyboardArrowUpIcon /> },
    closed: { title: "See More", icon: <KeyboardArrowDownIcon /> },
  },
  maxHeight = "none",
  ...props
}: CardContentProps) => {
  const ref = React.useRef(null);
  const [showExpand, setShowExpand] = React.useState<boolean>(false);
  const [expanded, setExpanded] = React.useState<boolean>(true);

  React.useEffect(() => {
    if (showExpand && expanded) {
      return;
    }
    if (ref.current.scrollHeight > maxHeight) {
      setShowExpand(true);
      setExpanded(false);
    }
  }, [maxHeight, children]);

  return (
    <BaseCardContent {...props} ref={ref}>
      <StyledCardContentContainer maxHeight={expanded ? "none" : maxHeight}>
        {children}
      </StyledCardContentContainer>
      {collapsible && showExpand && (
        <StyledCardActionArea padding={0} onClick={() => setExpanded(!expanded)}>
          <Grid container alignItems="center" justify="center">
            <Grid item>
              <Typography variant="body4" color="#3548D4">
                {expanded ? collapseAction?.open.title : collapseAction?.closed.title}
              </Typography>
            </Grid>
            <Grid item>
              <StyledExpandButton variant="neutral">
                {expanded ? collapseAction?.open.icon : collapseAction?.closed.icon}
              </StyledExpandButton>
            </Grid>
          </Grid>
        </StyledCardActionArea>
      )}
    </BaseCardContent>
  );
};

const StyledLandingCard = styled(Card)({
  border: "none",

  "& .header": {
    display: "inline-flex",
    marginBottom: "16px",
    fontWeight: "bold",
    fontSize: "12px",
    lineHeight: "36px",
    color: "rgba(13, 16, 48, 0.6)",
  },

  "& .header .icon .MuiAvatar-root": {
    height: "36px",
    width: "36px",
    marginRight: "8px",
    color: "rgba(13, 16, 48, 0.38)",
    backgroundColor: "rgba(13, 16, 48, 0.12)",
  },
});

export interface LandingCardProps extends Pick<CardActionAreaProps, "onClick"> {
  group: string;
  title: string;
  description: string;
}

export const LandingCard = ({ group, title, description, onClick, ...props }: LandingCardProps) => (
  <StyledLandingCard {...props}>
    <StyledCardActionArea onClick={onClick}>
      <CardContent padding={4}>
        <div className="header">
          <div className="icon">
            <Avatar>{group.charAt(0)}</Avatar>
          </div>
          <span>{group}</span>
        </div>
        <div>
          <Typography variant="h3">{title}</Typography>
          <Typography color="rgba(13, 16, 48, 0.6)" variant="body2">
            {description}
          </Typography>
        </div>
      </CardContent>
    </StyledCardActionArea>
  </StyledLandingCard>
);

export { Card, CardContent, CardHeader };
