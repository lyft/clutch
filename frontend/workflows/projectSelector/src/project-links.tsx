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
});

const StyledSubLink = styled.div({
  ...itemHoverStyle,
  paddingBottom: "8px",
  paddingTop: "8px",
  paddingLeft: "40px",
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
  padding: "7px 0px",
});

const StyledMultiLinkTitle = styled.span({
  fontWeight: "bold",
});

const StyledLinkBox = styled.div({
  borderRadius: "4px",
  width: "160px",
});

const StyledMultilinkImage = styled.div({
  display: "flex",
  padding: "8px",
});

const StyledMultilinkHeader = styled.div({
  display: "flex",
  alignItems: "center",
});

const StyledCenterImgSpan = styled.span({
  display: "flex",
  alignItems: "center",
  padding: "8px",
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
        <StyledCenterImgSpan>
          <img
            width={ICON_SIZE}
            height={ICON_SIZE}
            src={linkGroupImage}
            alt={validLinks[0].name ?? `Quick Link to ${validLinks[0].url}`}
          />
        </StyledCenterImgSpan>
        <StyledLinkTitle>{linkGroupName}</StyledLinkTitle>
      </Link>
    </StyledMenuItem>
  ) : (
    <div key={validLinks[0].url}>
      <StyledMultilinkHeader>
        <StyledMultilinkImage>
          <img width={ICON_SIZE} height={ICON_SIZE} src={linkGroupImage} alt={linkGroupName} />
        </StyledMultilinkImage>
        <StyledMultiLinkTitle>{linkGroupName}</StyledMultiLinkTitle>
      </StyledMultilinkHeader>
      <div>
        {validLinks.map(link => {
          return (
            link?.url && (
              <React.Fragment key={link.url}>
                <StyledSubLink>
                  <Link href={link.url}>
                    <Typography noWrap variant="body4">
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
});

interface QuickLinksPopperProps {
  /**
   * The linkgroups to render. They could be a mix of single
   * and multi-link groups.
   */
  linkGroups: IClutch.core.project.v1.ILinkGroup[];
  /**
   * A reference so that the popper knows to be attached to
   * the button.
   */
  anchorRef: React.RefObject<HTMLElement>;
  /** Whether the popper is open or not */
  open: boolean;
  /** A function that is called when closing / clicking away */
  onClose: () => void;
}

const QuickLinksPopper = ({ linkGroups, anchorRef, open, onClose }: QuickLinksPopperProps) => {
  return (
    <Popper open={open} anchorRef={anchorRef} onClickAway={onClose}>
      <ExpandedLinks linkGroups={linkGroups} />
    </Popper>
  );
};

interface ProjectLinksProps {
  /**
   * The linkgroups that will be rendered. They could be a mix
   * of single and multi-link groups.
   */
  linkGroups: IClutch.core.project.v1.ILinkGroup[];

  /**
   * A function that is called when the QuickLinksPopper opens.
   */
  onOpen: () => void;

  /**
   * A function that is called when the QuickLinksPopper
   * is closed.
   */
  onClose: () => void;

  /**
   * A boolean that denotes whether to render the button
   * that opens the quicklinks or not.
   */
  showOpenButton: boolean;
}

const ProjectLinks = ({ linkGroups, onOpen, onClose, showOpenButton }: ProjectLinksProps) => {
  const anchorRef = React.useRef(null);
  // The state is managed here for the popper because if it is hoisted up
  // to the parent that results in all the poppers being opened or closed
  // at the same time.
  const [open, setOpen] = React.useState(false);

  const onClosePopper = () => {
    setOpen(false);
    onClose();
  };

  const onOpenPopper = () => {
    setOpen(true);
    onOpen();
  };

  if (linkGroups.length === 0) {
    return null;
  }

  return (
    <StyledFlexEnd hidden={showOpenButton}>
      <StyledMoreVertIcon>
        <IconButton ref={anchorRef} onClick={onOpenPopper}>
          <MoreVertIcon />
          <QuickLinksPopper
            linkGroups={linkGroups}
            anchorRef={anchorRef}
            open={open}
            onClose={onClosePopper}
          />
        </IconButton>
      </StyledMoreVertIcon>
    </StyledFlexEnd>
  );
};

export default ProjectLinks;
