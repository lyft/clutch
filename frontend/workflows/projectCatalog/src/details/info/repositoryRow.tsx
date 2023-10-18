import React from "react";
import { Grid } from "@clutch-sh/core";
import type { IconProp } from "@fortawesome/fontawesome-svg-core";
import { faGithub } from "@fortawesome/free-brands-svg-icons";
import { faCode, faCodeBranch } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";

import { getRepositoryFromString, LinkText } from "../helpers";

interface ProjectPullRequests {
  number: number;
  url?: string;
}

const RepositoryRow = ({ repo }: { repo: string }) => {
  const [name, setName] = React.useState<string>("");
  const [url, setUrl] = React.useState<string>();
  const [icon, setIcon] = React.useState<IconProp>(faCode);

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const [requests, setRequests] = React.useState<ProjectPullRequests>();

  React.useEffect(() => {
    if (repo) {
      const { name: repoName, icon: repoIcon, url: repoUrl } = getRepositoryFromString(repo);
      setName(repoName);
      setIcon(repoIcon);
      setUrl(repoUrl);

      // TODO (jslaughter): fetch open PR's count
      // setRequests({
      //   number: 0,
      //   url: `https://${manager}/${project}/pulls`,
      // });
    }
  }, [repo]);

  return (
    <>
      <Grid item>
        <Grid container spacing={1}>
          <Grid item>{icon && <FontAwesomeIcon icon={icon} size="lg" />}</Grid>
          <Grid item>
            <LinkText text={name} link={url} />
          </Grid>
        </Grid>
      </Grid>
      {requests && (
        <Grid item>
          <Grid container spacing={1}>
            <Grid item>
              <FontAwesomeIcon icon={faCodeBranch} size="1x" />
            </Grid>
            <Grid item>
              <LinkText text={`${requests.number} open`} link={requests.url} />
            </Grid>
          </Grid>
        </Grid>
      )}
    </>
  );
};

export default RepositoryRow;
