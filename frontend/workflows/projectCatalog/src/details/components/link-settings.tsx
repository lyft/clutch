import React from "react";
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
  useParams,
} from "@clutch-sh/core";
import SettingsIcon from "@mui/icons-material/Settings";
import { isEmpty } from "lodash";

import type { ProjectConfigLink } from "../../types";

interface QuickLinksAndSettingsProps {
  linkGroups: QuickLinkGroup[];
  configLinks?: ProjectConfigLink[];
  showSettings: boolean;
}

const StyledPopperItem = styled(PopperItem)({
  "&&&": {
    height: "auto",
  },
  "& span.MuiTypography-root": {
    padding: "0",
  },
  "& a.MuiTypography-root": {
    padding: "4px 16px",
  },
});

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
      style={{
        padding: "8px 0px 0px 0px",
      }}
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
