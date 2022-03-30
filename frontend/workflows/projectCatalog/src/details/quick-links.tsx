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

interface QuickLinkContainerProps {
  key: string | null | undefined;
  name: string;
  children: React.ReactNode;
}

const ICON_SIZE = "32px";

const QuickLinkContainer = ({ key, name, children }: QuickLinkContainerProps) => {
  const container = (
    <Tooltip title={name}>
      <TooltipContainer>{children}</TooltipContainer>
    </Tooltip>
  );

  return (
    <Grid item key={key ?? ""}>
      {name ? container : children}
    </Grid>
  );
};

// If only a single link, then no popper is necessary
const QuickLink = ({ link, linkGroupName, linkGroupImage }: QuickLinkProps) =>
  link?.url ? (
    <QuickLinkContainer key={link.name} name={linkGroupName}>
      <Link href={link.url}>
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
  links: IClutch.core.project.v1.ILink[];
}
// Have a popper in the case of multiple links per group
const QuickLinkGroup = ({ linkGroupName, linkGroupImage, links }: QuickLinkGroupProps) => {
  const anchorRef = React.useRef(null);
  const [open, setOpen] = React.useState(false);
  const [validLinks, setValidLinks] = React.useState<IClutch.core.project.v1.ILink[]>([]);

  React.useEffect(() => {
    if (links) {
      setValidLinks(links.filter(link => link?.url && link.url.length > 0));
    }
  }, [links]);

  return (
    <QuickLinkContainer key={linkGroupName} name={linkGroupName}>
      <button
        type="button"
        style={{ padding: 0, background: "transparent", border: "0", cursor: "pointer" }}
        ref={anchorRef}
        onClick={() => setOpen(true)}
      >
        <img width={ICON_SIZE} height={ICON_SIZE} src={linkGroupImage} alt={linkGroupName} />
      </button>
      <Popper open={open} anchorRef={anchorRef} onClickAway={() => setOpen(false)}>
        {validLinks.map(link => (
          <PopperItem key={link.name}>
            {link?.url && (
              <Link href={link.url}>
                <Typography color="inherit" variant="body4">
                  {link.name}
                </Typography>
              </Link>
            )}
          </PopperItem>
        ))}
      </Popper>
    </QuickLinkContainer>
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
