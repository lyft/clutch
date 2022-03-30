import React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";
import {
  Card,
  Grid,
  Link,
  Popper,
  PopperItem,
  Tooltip,
  TooltipContainer,
  Typography,
} from "@clutch-sh/core";

interface LinkGroupProps {
  linkGroupName: string;
  linkGroupImage: string;
}

interface QuickLinkProps extends LinkGroupProps {
  link: IClutch.core.project.v1.ILink;
}

interface QuickLinkTooltipProps {
  key: string;
  name: string;
  children: React.ReactNode;
}

const ICON_SIZE = "32px";

const QuickLinkTooltip = ({ key, name, children }: QuickLinkTooltipProps) => (
  <Grid item key={key}>
    <Tooltip title={name}>
      <TooltipContainer>{children}</TooltipContainer>
    </Tooltip>
  </Grid>
);

// If only a single link, then no popper is necessary
const QuickLink = ({ link, linkGroupName, linkGroupImage }: QuickLinkProps) => {
  const container = (
    <Link href={link.url ?? undefined}>
      <img
        width={ICON_SIZE}
        height={ICON_SIZE}
        src={linkGroupImage}
        alt={link.name ?? "Quick Link"}
      />
    </Link>
  );
  return linkGroupName ? (
    <QuickLinkTooltip key={link.name ?? ""} name={linkGroupName}>
      {container}
    </QuickLinkTooltip>
  ) : (
    container
  );
};

interface QuickLinkGroupProps extends LinkGroupProps {
  links: IClutch.core.project.v1.ILink[];
}
// Have a popper in the case of multiple links per group
const QuickLinkGroup = ({ linkGroupName, linkGroupImage, links }: QuickLinkGroupProps) => {
  const anchorRef = React.useRef(null);
  const [open, setOpen] = React.useState(false);

  const container = (
    <>
      <button
        type="button"
        style={{ padding: 0, background: "transparent", border: "0", cursor: "pointer" }}
        ref={anchorRef}
        onClick={() => setOpen(true)}
      >
        <img width={ICON_SIZE} height={ICON_SIZE} src={linkGroupImage} alt={linkGroupName} />
      </button>
      <Popper open={open} anchorRef={anchorRef} onClickAway={() => setOpen(false)}>
        {links.map(link => (
          <PopperItem key={link.name}>
            <Link href={link.url ?? undefined}>
              <Typography color="inherit" variant="body4">
                {link.name}
              </Typography>
            </Link>
          </PopperItem>
        ))}
      </Popper>
    </>
  );

  return linkGroupName ? (
    <QuickLinkTooltip key={linkGroupName} name={linkGroupName}>
      {container}
    </QuickLinkTooltip>
  ) : (
    container
  );
};
export interface QuickLinksProps {
  linkGroups: IClutch.core.project.v1.ILinkGroup[];
}

const QuickLinksCard = ({ linkGroups }: QuickLinksProps) => (
  <Card>
    <Grid
      container
      item
      direction="column"
      alignItems="center"
      spacing={1}
      style={{ padding: "8px" }}
    >
      {(linkGroups || []).map(linkGroup => {
        if (linkGroup.links?.length === 1) {
          return (
            <QuickLink
              link={linkGroup.links[0]}
              linkGroupName={linkGroup.name ?? ""}
              linkGroupImage={linkGroup.imagePath ?? ""}
            />
          );
        }
        return (
          <QuickLinkGroup
            linkGroupName={linkGroup.name ?? ""}
            linkGroupImage={linkGroup.imagePath ?? ""}
            links={linkGroup?.links ?? []}
          />
        );
      })}
    </Grid>
  </Card>
);

export default QuickLinksCard;
