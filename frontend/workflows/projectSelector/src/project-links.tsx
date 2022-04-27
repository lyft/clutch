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

const ICON_SIZE = "16px";

const StyledMenuItem = styled.div({
  display: "flex",
  alignItems: "center",
  "&:hover": {
    backgroundColor: "rgba(13, 16, 48, 0.05)",
  },
});

const StyledMoreVertIcon = styled.span({
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
  // In the case where there is only a singe link in the group, we make the title clickable.
  // In the case where there are multiple links, the title is not clickable and has different styling.
  return links.length === 1 ? (
    <StyledMenuItem style={{ padding: "8px" }} key={links[0].url}>
      <Link href={links[0]?.url ?? ""}>
        <img
          width={ICON_SIZE}
          height={ICON_SIZE}
          src={linkGroupImage}
          alt={links[0].name ?? `Quick Link to ${links[0].url}`}
        />
        <StyledLinkTitle style={{ paddingLeft: "6px" }}>{linkGroupName}</StyledLinkTitle>
      </Link>
    </StyledMenuItem>
  ) : (
    <div>
      <div style={{ display: "flex", alignItems: "center" }}>
        <div style={{ paddingLeft: "6px", paddingRight: "6px", paddingTop: "6px" }}>
          <img width={ICON_SIZE} height={ICON_SIZE} src={linkGroupImage} alt={linkGroupName} />
        </div>
        <StyledLinkTitle>{linkGroupName}</StyledLinkTitle>
      </div>
      <div>
        {validLinks.map(link => {
          return (
            link?.url && (
              <React.Fragment key={link.url}>
                <StyledMenuItem
                  style={{ paddingBottom: "4px", paddingRight: "4px", paddingLeft: "46px" }}
                >
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

const StyledFlexEnd = styled.div({
  justifyContent: "right",
  flexDirection: "row-reverse",
});

const QuickLinksPopper = ({ linkGroups, anchorRef, open, setOpen, setKeyWithQLinksOpen }) => {
  return (
    <Popper
      open={open}
      anchorRef={anchorRef}
      onClickAway={() => {
        setOpen(false);
        setKeyWithQLinksOpen("");
      }}
    >
      <ExpandedLinks linkGroups={linkGroups} />
    </Popper>
  );
};

const ProjectLinks = ({ linkGroups, name }) => {
  const anchorRef = React.useRef(null);
  const [open, setOpen] = React.useState(false);
  const [keyWithQLinksOpen, setKeyWithQLinksOpen] = React.useState("");

  return (
    <StyledFlexEnd hidden={name !== keyWithQLinksOpen}>
      <StyledMoreVertIcon>
        <IconButton
          ref={anchorRef}
          onClick={() => {
            setOpen(true);
            setKeyWithQLinksOpen(name);
          }}
        >
          <MoreVertIcon />
          <QuickLinksPopper
            linkGroups={linkGroups}
            anchorRef={anchorRef}
            open={open}
            setOpen={setOpen}
            setKeyWithQLinksOpen={setKeyWithQLinksOpen}
          />
        </IconButton>
      </StyledMoreVertIcon>
    </StyledFlexEnd>
  );
};

export default ProjectLinks;
