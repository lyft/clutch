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

const itemHoverStyle = {
  display: "flex",
  alignItems: "center",
  "&:hover": {
    backgroundColor: "rgba(13, 16, 48, 0.05)",
  },
};

const StyledMenuItem = styled.div({
  ...itemHoverStyle,
  padding: "8px",
});

const StyledSubLink = styled.div({
  ...itemHoverStyle,
  paddingBottom: "4px",
  paddingRight: "4px",
  paddingLeft: "46px",
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

const StyledLinkBox = styled.div({
  borderRadius: "4px",
});

const StyledMultilinkImage = styled.div({
  paddingLeft: "6px",
  paddingRight: "6px",
  paddingTop: "6px",
});

const StyledMultilinkTitle = styled.div({
  display: "flex",
  alignItems: "center",
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

  if (validLinks.length === 0) {
    return null;
  }

  // In the case where there is only a singe link in the group, we make the title clickable.
  // In the case where there are multiple links, the title is not clickable and has different styling.
  return validLinks.length === 1 ? (
    <StyledMenuItem key={validLinks[0].url}>
      <Link href={validLinks[0]?.url ?? ""}>
        <img
          width={ICON_SIZE}
          height={ICON_SIZE}
          src={linkGroupImage}
          alt={validLinks[0].name ?? `Quick Link to ${validLinks[0].url}`}
        />
        <StyledLinkTitle style={{ paddingLeft: "6px" }}>{linkGroupName}</StyledLinkTitle>
      </Link>
    </StyledMenuItem>
  ) : (
    <div key={validLinks[0].url}>
      <StyledMultilinkTitle>
        <StyledMultilinkImage>
          <img width={ICON_SIZE} height={ICON_SIZE} src={linkGroupImage} alt={linkGroupName} />
        </StyledMultilinkImage>
        <StyledLinkTitle>{linkGroupName}</StyledLinkTitle>
      </StyledMultilinkTitle>
      <div>
        {validLinks.map(link => {
          return (
            link?.url && (
              <React.Fragment key={link.url}>
                <StyledSubLink>
                  <Link href={link.url}>
                    <Typography color="inherit" variant="body4">
                      {link.name}
                    </Typography>
                  </Link>
                </StyledSubLink>
              </React.Fragment>
            )
          );
        })}
      </div>
    </div>
  );
};

interface ExpandedLinksProps {
  linkGroups: IClutch.core.project.v1.ILinkGroup[];
}

const ExpandedLinks = ({ linkGroups }: ExpandedLinksProps) => (
  <StyledLinkBox>
    {(linkGroups || []).map(linkGroup => {
      return (
        <QuickLinkGroup
          linkGroupName={linkGroup.name ?? ""}
          linkGroupImage={linkGroup.imagePath ?? ""}
          links={linkGroup?.links ?? []}
        />
      );
    })}
  </StyledLinkBox>
);

const StyledFlexEnd = styled.div({
  justifyContent: "right",
  flexDirection: "row-reverse",
});

interface QuickLinksPopperProps {
  linkGroups: IClutch.core.project.v1.ILinkGroup[];
  anchorRef: React.RefObject<HTMLElement>;
  open: boolean;
  setOpen: (open: boolean) => void;
  setKeyWithQLinksOpen: (key: string) => void;
}

const QuickLinksPopper = ({
  linkGroups,
  anchorRef,
  open,
  setOpen,
  setKeyWithQLinksOpen,
}: QuickLinksPopperProps) => {
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

interface ProjectLinksProps {
  linkGroups: IClutch.core.project.v1.ILinkGroup[];
  name: string[];
}

const ProjectLinks = ({ linkGroups, name }: ProjectLinksProps) => {
  const anchorRef = React.useRef(null);
  const [open, setOpen] = React.useState(false);
  const [keyWithQLinksOpen, setKeyWithQLinksOpen] = React.useState("");

  return (
    <StyledFlexEnd hidden={name?.[0] !== keyWithQLinksOpen}>
      <StyledMoreVertIcon>
        <IconButton
          ref={anchorRef}
          onClick={() => {
            setOpen(true);
            setKeyWithQLinksOpen(name?.[0]);
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
