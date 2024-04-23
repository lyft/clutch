import React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import { Badge, BadgeProps, SvgIcon } from "@mui/material";

import { IconButton } from "./button";
import { Card } from "./card";
import { Tooltip, TooltipContainer } from "./Feedback";
import Grid from "./grid";
import { Link } from "./link";
import { Popper, PopperItem } from "./popper";
import styled from "./styled";
import { Typography } from "./typography";

const StyledGrid = styled(Grid)({
  padding: "10px",
  margin: "-4px",
});

const StyledQLGrid = styled(Grid)({
  padding: "8px",
});

const StyledButton = styled("button")({
  padding: 0,
  background: "transparent",
  border: "0",
  cursor: "pointer",
  display: "flex",
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

const StyledBadge = styled(Badge)({
  ".MuiBadge-anchorOriginTopRightCircular": {
    top: "23%",
    right: "23%",
  },
  ".MuiBadge-dot": {
    height: "10px",
    minWidth: "10px",
    borderRadius: "50%",
  },
});

interface LinkGroupProps {
  linkGroupName: string;
  linkGroupImage: string;
}

export interface QLink extends IClutch.core.project.v1.ILink {
  trackingId?: string;
  icon?: React.ElementType;
}

interface QuickLinkProps extends LinkGroupProps {
  link: QLink;
}

interface QuickLinkContainerProps {
  keyProp: string | null | undefined;
  name: string;
  children: React.ReactNode;
  popperOpen?: boolean;
}

const ICON_SIZE = "32px";

const QuickLinkContainer = ({ keyProp, name, children, popperOpen }: QuickLinkContainerProps) => {
  const [tooltipOpen, setTooltipOpen] = React.useState<boolean>(false);

  React.useEffect(() => {
    if (popperOpen) {
      setTooltipOpen(false);
    }
  }, [popperOpen]);

  const container = (
    <Tooltip title={name} open={tooltipOpen}>
      <TooltipContainer
        onMouseEnter={() => !popperOpen && setTooltipOpen(true)}
        onMouseLeave={() => !popperOpen && setTooltipOpen(false)}
      >
        {children}
      </TooltipContainer>
    </Tooltip>
  );

  return (
    <StyledQLGrid item key={keyProp ?? ""}>
      {name ? container : children}
    </StyledQLGrid>
  );
};

const QuickLinkWrapper = ({ linkGroup, children }) => (
  <StyledBadge
    key={`quicklink-${linkGroup.name}`}
    badgeContent={linkGroup.badge?.content ?? null}
    color={linkGroup.badge?.color ?? "default"}
    overlap="circular"
    variant={linkGroup.badge?.content.trim() ? "standard" : "dot"}
  >
    {children}
  </StyledBadge>
);

// If only a single link, then no popper is necessary
const QuickLink = ({ link, linkGroupName, linkGroupImage }: QuickLinkProps) =>
  link?.url ? (
    <QuickLinkContainer keyProp={link.name} name={linkGroupName}>
      <Link href={link.url} data-tracking-action={link.trackingId}>
        <img
          width={ICON_SIZE}
          height={ICON_SIZE}
          src={linkGroupImage}
          alt={link.name ?? `Quick Link to ${link.url}`}
        />
      </Link>
    </QuickLinkContainer>
  ) : null;

interface QuickLinkGroupProps extends LinkGroupProps {
  links: QLink[];
}
// Have a popper in the case of multiple links per group
const QuickLinkGroup = ({ linkGroupName, linkGroupImage, links }: QuickLinkGroupProps) => {
  const anchorRef = React.useRef(null);
  const [open, setOpen] = React.useState(false);
  const [validLinks, setValidLinks] = React.useState<QLink[]>([]);

  React.useEffect(() => {
    if (links) {
      setValidLinks(links.filter(link => link?.url && link.url.length > 0));
    }
  }, [links]);

  return (
    <QuickLinkContainer keyProp={linkGroupName} name={linkGroupName} popperOpen={open}>
      <StyledButton type="button" ref={anchorRef} onClick={() => setOpen(true)}>
        <img width={ICON_SIZE} height={ICON_SIZE} src={linkGroupImage} alt={linkGroupName} />
      </StyledButton>
      <Popper
        open={open}
        anchorRef={anchorRef}
        onClickAway={() => setOpen(false)}
        placement="bottom-end"
      >
        {validLinks.map(link => (
          <StyledPopperItem key={link.name}>
            {link?.url && (
              <Link href={link.url}>
                <Grid container alignItems="center" gap={1}>
                  {link?.icon && <SvgIcon component={link.icon} fontSize="small" />}
                  <Typography color="inherit" variant="body4" noWrap>
                    {link.name}
                  </Typography>
                </Grid>
              </Link>
            )}
          </StyledPopperItem>
        ))}
      </Popper>
    </QuickLinkContainer>
  );
};

export interface LinkGroup extends IClutch.core.project.v1.ILinkGroup {
  links?: QLink[];
  badge?: {
    color: BadgeProps["color"];
    content: string;
  };
}

export interface QuickLinksProps {
  linkGroups: LinkGroup[];
  maxLinks?: number;
}

// TODO(smonero): Wasn't sure if I should make an interface for this or just reuse
// or not make one at all since its so simple
interface SlicedLinkGroupProps {
  slicedLinkGroups: LinkGroup[];
}

const SlicedLinkGroup = ({ slicedLinkGroups }: SlicedLinkGroupProps) => {
  return (
    <>
      {(slicedLinkGroups || []).map(linkGroup => {
        return (
          <QuickLinkWrapper linkGroup={linkGroup} key={linkGroup.name}>
            {linkGroup.links?.length === 1 ? (
              <QuickLink
                key={`quicklink-${linkGroup.name}`}
                link={linkGroup.links[0]}
                linkGroupName={linkGroup.name ?? ""}
                linkGroupImage={linkGroup.imagePath ?? ""}
              />
            ) : (
              <QuickLinkGroup
                key={`quicklink-${linkGroup.name}`}
                linkGroupName={linkGroup.name ?? " "}
                linkGroupImage={linkGroup.imagePath ?? ""}
                links={linkGroup?.links ?? []}
              />
            )}
          </QuickLinkWrapper>
        );
      })}
    </>
  );
};

const QuickLinksCard = ({ linkGroups, maxLinks = 5 }: QuickLinksProps) => {
  const anchorRef = React.useRef(null);
  const [open, setOpen] = React.useState(false);

  const filteredLinkGroups = linkGroups.filter(
    lg => lg.links?.length > 0 && lg.name && lg.imagePath
  );

  // Show only the first five quick links, and put the rest in
  // an overflow popper
  const firstFive = filteredLinkGroups.slice(0, maxLinks);
  const overflow = filteredLinkGroups.slice(maxLinks);

  return (
    <Card>
      <StyledGrid
        container
        item
        direction="row"
        alignItems="center"
        spacing={1}
        justifyContent="space-around"
        flexWrap="nowrap"
      >
        <SlicedLinkGroup slicedLinkGroups={firstFive} />
        {overflow.length > 0 && (
          <>
            <IconButton
              size="small"
              variant="neutral"
              ref={anchorRef}
              onClick={() => setOpen(true)}
            >
              <ExpandMoreIcon />
            </IconButton>
            <Popper
              open={open}
              anchorRef={anchorRef}
              onClickAway={() => setOpen(false)}
              placement="bottom-end"
            >
              <StyledQLGrid direction="row" container>
                <SlicedLinkGroup slicedLinkGroups={overflow} />
              </StyledQLGrid>
            </Popper>
          </>
        )}
      </StyledGrid>
    </Card>
  );
};
export default QuickLinksCard;
