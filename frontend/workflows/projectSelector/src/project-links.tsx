import React from "react";
import type { clutch as IClutch } from "@clutch-sh/api";
import { Link, Popper, Tooltip, TooltipContainer, Typography } from "@clutch-sh/core";
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

interface QuickLinkContainerProps {
  key: string | null | undefined;
  name: string;
  children: React.ReactNode;
}

const ICON_SIZE = "16px";

const StyledMenuItem = styled.div({
  display: "flex",
  alignItems: "center",
  "&:hover": {
    backgroundColor: "rgba(13, 16, 48, 0.05)",
  },
});

const QuickLinkContainer = ({ key, name, children }: QuickLinkContainerProps) => {
  const container = (
    <Tooltip title={name}>
      <TooltipContainer>{children}</TooltipContainer>
    </Tooltip>
  );

  return <div key={key}>{name ? container : children}</div>;
};

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
  padding: "8px",
  fontWeight: "bold",
});

const QuickLink = ({ link, linkGroupName, linkGroupImage }: QuickLinkProps) =>
  link?.url ? (
    <StyledMenuItem style={{ padding: "8px" }}>
      <Link href={link.url}>
        <QuickLinkContainer key={link.name} name={linkGroupName}>
          <img
            width={ICON_SIZE}
            height={ICON_SIZE}
            src={linkGroupImage}
            alt={link.name ?? `Quick Link to ${link.url}`}
          />
          <StyledLinkTitle>{linkGroupName}</StyledLinkTitle>
        </QuickLinkContainer>
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
    <QuickLinkContainer key={linkGroupName} name={linkGroupName}>
      <StyledLinkTitle>
        <img width={ICON_SIZE} height={ICON_SIZE} src={linkGroupImage} alt={linkGroupName} />
      </StyledLinkTitle>
      <StyledLinkTitle style={{ padding: "0px" }}>{linkGroupName}</StyledLinkTitle>
      <div>
        {validLinks.map(link => {
          return (
            link?.url && (
              <StyledMenuItem style={{ padding: "8px", paddingLeft: "40px" }}>
                <Link href={link.url}>
                  <Typography color="inherit" variant="body4">
                    {link.name}
                  </Typography>
                </Link>
              </StyledMenuItem>
            )
          );
        })}
      </div>
    </QuickLinkContainer>
  );
};

const ExpandedLinks = ({ linkGroups }) => (
  <div>
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
