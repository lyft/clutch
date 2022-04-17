import React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";
import { Link, Popper, Typography } from "@clutch-sh/core";
import styled from "@emotion/styled";
import IconButton from "@material-ui/core/IconButton";
import MoreVertIcon from "@material-ui/icons/MoreVert";

interface LinkGroupProps {
  linkGroupName: string;
  linkGroupImage: string;
}

interface QuickLinkProps extends LinkGroupProps {
  link: IClutch.core.project.v1.ILink;
}

const ICON_SIZE = "22px";

const StyledMenuItem = styled.div({
  display: "flex",
  alignItems: "center",
  "&:hover": {
    backgroundColor: "rgba(13, 16, 48, 0.05)",
  },
});

const StyledMoreVertIcon = styled.span({
  width: "12px",
  ".MuiIconButton-root": {
    padding: "6px",
    color: "rgba(13, 16, 48, 0.38)",
  },
  ".MuiIconButton-root:hover": {
    backgroundColor: "rgb(245, 246, 253)",
  },
  ".MuiIconButton-root:active": {
    backgroundColor: "rgba(0,0,0, 0.1)",
  },
});

const StyledLinkTitle = styled.span({
  fontWeight: "bold",
});

const QuickLink = ({ link, linkGroupName, linkGroupImage }: QuickLinkProps) =>
  link?.url ? (
    <StyledMenuItem style={{ padding: "8px" }} key={link.url}>
      <Link href={link.url}>
        <img
          width={ICON_SIZE}
          height={ICON_SIZE}
          src={linkGroupImage}
          alt={link.name ?? `Quick Link to ${link.url}`}
        />
        <StyledLinkTitle style={{ paddingLeft: "6px" }}>{linkGroupName}</StyledLinkTitle>
      </Link>
    </StyledMenuItem>
  ) : null;

interface QuickLinkGroupProps extends LinkGroupProps {
  links: IClutch.core.project.v1.ILink[];
}

const QuickLinkGroup = ({ linkGroupName, linkGroupImage, links }: QuickLinkGroupProps) => {
  const [validLinks, setValidLinks] = React.useState<IClutch.core.project.v1.ILink[]>([]);

  React.useEffect(() => {
    if (links) {
      setValidLinks(links.filter(link => link?.url && link.url.length > 0));
    }
  }, [links]);
  return (
    <div key={linkGroupName}>
      <div style={{ display: "flex", alignItems: "center" }}>
        <div style={{ padding: "6px" }}>
          <img width={ICON_SIZE} height={ICON_SIZE} src={linkGroupImage} alt={linkGroupName} />
        </div>
        <StyledLinkTitle>{linkGroupName}</StyledLinkTitle>
      </div>
      <div>
        {validLinks.map(link => {
          return (
            link?.url && (
              <React.Fragment key={link.url}>
                <StyledMenuItem style={{ padding: "6px", paddingLeft: "46px" }}>
                  <Link href={link.url}>
                    <Typography color="inherit" variant="body4">
                      {link.name}
                    </Typography>
                  </Link>
                </StyledMenuItem>
              </React.Fragment>
            )
          );
        })}
      </div>
    </div>
  );
};

const ExpandedLinks = ({ linkGroups }) => (
  <div style={{ borderRadius: "4px" }}>
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
  </div>
);

const ProjectLinks = ({ linkGroups }) => {
  const anchorRef = React.useRef(null);
  const [open, setOpen] = React.useState(false);

  if (!linkGroups) {
    return null;
  }

  return (
    <StyledMoreVertIcon>
      <IconButton ref={anchorRef} onClick={() => setOpen(true)}>
        <MoreVertIcon />
        <Popper open={open} anchorRef={anchorRef} onClickAway={() => setOpen(false)}>
          <ExpandedLinks linkGroups={linkGroups} />
        </Popper>
      </IconButton>
    </StyledMoreVertIcon>
  );
};

export default ProjectLinks;
