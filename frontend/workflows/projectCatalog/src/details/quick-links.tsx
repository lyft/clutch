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

// If only a single link, then no popper is necessary
const QuickLink = ({ link, linkGroupName, linkGroupImage }) => {
  <Grid item key={link.name}>
    <Tooltip title={linkGroupName}>
      <TooltipContainer>
        <Link href={link.url}>
          <img width="32px" height="32px" src={linkGroupImage} alt={link.name} />
        </Link>
      </TooltipContainer>
    </Tooltip>
  </Grid>;
};

// Have a popper in the case of multiple links per group
const QuickLinkGroup = ({ linkGroupName, linkGroupImage, links }) => {
  const anchorRef = React.useRef(null);
  const [open, setOpen] = React.useState(false);

  return (
    <Grid item key={linkGroupName}>
      <Tooltip title={linkGroupName}>
        <TooltipContainer>
          {/* eslint-disable */}
          <img
            width="32px"
            height="32px"
            src={linkGroupImage}
            alt={linkGroupName}
            onClick={() => setOpen(true)}
            ref={anchorRef}
          />
          <Popper open={open} anchorRef={anchorRef} onClickAway={() => setOpen(false)}>
            {links.map(link => {
              return (
                <PopperItem key={link.name}>
                  <Link href={link.url}>
                    <Typography color="inherit" variant="body4">
                      {link.name}
                    </Typography>
                  </Link>
                </PopperItem>
              );
            })}
          </Popper>
          </TooltipContainer>
      </Tooltip>
    </Grid>
  );
};
export interface QuickLinksCardInput {
  linkGroups: IClutch.core.project.v1.ILinkGroup[];
}

const QuickLinksCard = ({ linkGroups }: QuickLinksCardInput) => {
  return (
    <Card>
      <Grid
        container
        item
        direction="column"
        alignItems="center"
        spacing={1}
        style={{ padding: "7px 0" }}
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
};

export default QuickLinksCard;
