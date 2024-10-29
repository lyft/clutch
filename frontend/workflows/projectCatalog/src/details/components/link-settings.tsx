import React from "react";
import { useParams } from "react-router-dom";
import {
  checkFeatureEnabled,
  Grid,
  IconButton,
  Link,
  Popper,
  PopperItem,
  QuickLinkGroup,
  QuickLinksCard,
  styled,
  Typography,
} from "@clutch-sh/core";
import SettingsIcon from "@mui/icons-material/Settings";
import type { Theme } from "@mui/material";
import { isEmpty } from "lodash";

import type { ProjectConfigLink } from "../../types";

interface QuickLinksAndSettingsProps {
  linkGroups: QuickLinkGroup[];
  configLinks?: ProjectConfigLink[];
  showSettings: boolean;
}

const StyledPopperItem = styled(PopperItem)(({ theme }: { theme: Theme }) => ({
  "&&&": {
    height: "auto",
  },
  "& span.MuiTypography-root": {
    padding: theme.spacing(theme.clutch.spacing.none),
  },
  "& a.MuiTypography-root": {
    padding: theme.spacing(theme.clutch.spacing.xs, theme.clutch.spacing.base),
  },
}));

const QuickLinksAndSettings = ({
  linkGroups,
  configLinks = [],
  showSettings,
}: QuickLinksAndSettingsProps) => {
  const { projectId } = useParams();
  const anchorRef = React.useRef(null);
  const [open, setOpen] = React.useState(false);
  const [links, setLinks] = React.useState<ProjectConfigLink[]>(configLinks);

  React.useEffect(() => {
    const projectConfigFlag = checkFeatureEnabled({ feature: "projectCatalogSettings" });
    if (projectConfigFlag) {
      setLinks([
        {
          title: "Project Configuration",
          path: `/catalog/${projectId}/config`,
          icon: <SettingsIcon fontSize="small" />,
        },
        ...links,
      ]);
    }
  }, []);

  return (
    <Grid
      container
      direction="row"
      alignItems="center"
      justifyContent="flex-end"
      spacing={1}
      paddingTop={1}
    >
      {!isEmpty(linkGroups) && (
        <Grid item>
          <QuickLinksCard linkGroups={linkGroups} />
        </Grid>
      )}
      {showSettings && links && links.length > 0 && (
        <Grid item>
          <IconButton ref={anchorRef} onClick={() => setOpen(o => !o)} size="medium">
            <SettingsIcon />
          </IconButton>
          <Popper
            open={open}
            anchorRef={anchorRef}
            onClickAway={() => setOpen(false)}
            placement="bottom-end"
          >
            {links.map(link => (
              <StyledPopperItem key={link.title}>
                <Link href={link.path} target="_self">
                  <Grid container gap={0.5}>
                    {link.icon && <Grid item>{link.icon}</Grid>}
                    <Grid item>
                      <Typography variant="body2" color="inherit">
                        {link.title}
                      </Typography>
                    </Grid>
                  </Grid>
                </Link>
              </StyledPopperItem>
            ))}
          </Popper>
        </Grid>
      )}
    </Grid>
  );
};

export default QuickLinksAndSettings;
