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

const ICON_SIZE = "32px";
// If only a single link, then no popper is necessary
const QuickLink = ({ link, linkGroupName, linkGroupImage }) => (
  <Grid item key={link.name}>
    <Tooltip title={linkGroupName}>
      <TooltipContainer>
        <Link href={link.url}>
          <img width={ICON_SIZE} height={ICON_SIZE} src={linkGroupImage} alt={link.name} />
        </Link>
      </TooltipContainer>
    </Tooltip>
  </Grid>
);

// Have a popper in the case of multiple links per group
const QuickLinkGroup = ({ linkGroupName, linkGroupImage, links }) => {
  const anchorRef = React.useRef(null);
  const [open, setOpen] = React.useState(false);

  return (
    <Grid item key={linkGroupName}>
      <Tooltip title={linkGroupName}>
        <TooltipContainer>
          <button
            style={{ background: `url${linkGroupImage}` }}
            ref={anchorRef}
            onClick={() => setOpen(true)}
            type="button"
            aria-label="expand link group"
          />
          <Popper open={open} anchorRef={anchorRef} onClickAway={() => setOpen(false)}>
            {links.map(link => (
              <PopperItem key={link.name}>
                <Link href={link.url}>
                  <Typography color="inherit" variant="body4">
                    {link.name}
                  </Typography>
                </Link>
              </PopperItem>
            ))}
          </Popper>
        </TooltipContainer>
      </Tooltip>
    </Grid>
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
      {linkGroups?.map(linkGroup => {
        if (linkGroup.links?.length === 1) {
          return (
            <QuickLink
              link={linkGroup.links[0]}
              linkGroupName={linkGroup.name}
              linkGroupImage={linkGroup.imagePath}
            />
          );
        }
        return (
          <QuickLinkGroup
            linkGroupName={linkGroup.name}
            linkGroupImage={linkGroup.imagePath}
            links={linkGroup.links}
          />
        );
      })}
    </Grid>
  </Card>
);

export default QuickLinksCard;
