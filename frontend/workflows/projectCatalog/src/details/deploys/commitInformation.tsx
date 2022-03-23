import React from "react";
import { Grid, Link, styled, Tooltip, Typography } from "@clutch-sh/core";
import { faCodeCommit } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { Avatar } from "@material-ui/core";
import AvatarGroup from "@material-ui/lab/AvatarGroup";
import { uniqBy } from "lodash";

import EnvironmentIcon from "./environmentIcon";
import type { Commit, CommitInfo } from "./types";

const unknownUser = "Unknown User";
const githubBaseURL = "https://github.com/";

const Icon = styled("div")({
  height: "20px",
  width: "20px",
});

const StyledAvatarGroup = styled(AvatarGroup)({
  ".MuiAvatarGroup-avatar": {
    height: "20px",
    width: "20px",
    border: "unset",
    fontSize: "9px",
  },
});

const StyledAvatar = styled(Avatar)({
  height: "20px",
  width: "20px",
  border: "unset",
  fontSize: "12px",
});

const StyledCommitContainer = styled("span")({
  width: "100%",
});

const StyledCommitRange = styled("div")({
  display: "flex",
  alignItems: "center",
  flexWrap: "wrap",
  "> :first-of-type": {
    paddingRight: "8px",
  },
});

// In many cases we only need the first 7 chars of a SHA
const shortSha = (s?: string | null) => {
  return s ? s.substr(0, 7) : "";
};

const CommitMessageAndNumberDiv = styled("div")({
  display: "flex",
});

const CommitMessageTruncateDiv = styled("div")({
  textOverflow: "ellipsis",
  overflow: "hidden",
  whiteSpace: "nowrap",
});

const SpaceIconDiv = styled("div")({
  marginLeft: "4px",
});

const commitMessage = (commits: Commit[]) => {
  if (!commits || commits.length < 1) {
    return "No commits found";
  }
  const lastCommitMessage = commits[commits.length - 1]?.message;
  // Extract the message and the commit number from the "message" field
  // If there is no regex match (e.g. message is empty) default to the whole message.
  const commitMessageRegex = /^(.*?)\((#\d+)\)/;
  const matched = lastCommitMessage?.match(commitMessageRegex);

  if (!matched || matched.length < 3) {
    return <CommitMessageTruncateDiv>{lastCommitMessage}</CommitMessageTruncateDiv>;
  }

  const justMessage = matched[1];
  const justNumber = matched[2];
  return (
    <CommitMessageAndNumberDiv>
      <CommitMessageTruncateDiv>{justMessage}</CommitMessageTruncateDiv>
      <SpaceIconDiv>({justNumber})</SpaceIconDiv>
    </CommitMessageAndNumberDiv>
  );
};

const CommitInformation = ({ repositoryName, commits, baseRef, environment }: CommitInfo) => {
  const commitCount = commits.length;
  const firstRef = shortSha(commits?.[0]?.ref);
  const lastRef = shortSha(commits?.[commits.length - 1]?.ref);

  const linkPath = commitCount === 1 ? `commit/${firstRef}` : `compare/${baseRef}...${lastRef}`;
  // Filter out commits missing authors
  const validCommits = commits.filter(c => c.author);
  // Deduplicate the commit authors
  const uniqueCommits = uniqBy(validCommits, "author.username");

  // If the authors email is a github noreply, we can parse it to find their github username.
  // Which can be used to lookup their avatar reliably.
  //
  // There are two different formats:
  // EG:
  //    example@users.noreply.github.com
  //    39421794+example@users.noreply.github.com
  const ghNoReplyRxp = /([0-9]*[+]?)(.*)(@users.noreply.github.com)/g;
  const authors = uniqueCommits.map(c => {
    let username = c.author?.username;
    const emailParts = ghNoReplyRxp.exec(c.author?.email || "");
    // if the domain is from github noreply, get the username which is always in group 2 of the regex.
    if (emailParts?.length === 4 && emailParts[3] === "@users.noreply.github.com") {
      const { 2: uname } = emailParts;
      username = uname;
    }

    return {
      username,
      email: c.author?.email,
    };
  });

  return (
    <StyledCommitContainer>
      <Typography variant="input">{commitMessage(commits)}</Typography>
      {commits && commits.length > 0 && (
        <Grid container direction="row" alignItems="center" spacing={2}>
          {environment && (
            <Grid item>
              <Tooltip title={environment}>
                <Icon>{EnvironmentIcon(environment)}</Icon>
              </Tooltip>
            </Grid>
          )}
          <Grid item>
            <StyledCommitRange>
              <Link href={`http://github.com/${repositoryName}/${linkPath}`}>
                <Typography variant="body4" color="#3548D4">
                  {lastRef}
                </Typography>
              </Link>
              <FontAwesomeIcon icon={faCodeCommit} />
              <Typography variant="body4">{commitCount}</Typography>
              <StyledAvatarGroup max={4} spacing={4}>
                {authors.length !== 0 ? (
                  authors.map((u, idx) => (
                    <Tooltip title={u.username || "unknown"} key={u.username || idx}>
                      <StyledAvatar
                        alt={u?.username || "username"}
                        src={`${githubBaseURL + u?.username?.replaceAll(" ", "")}.png`}
                      />
                    </Tooltip>
                  ))
                ) : (
                  <Tooltip title={unknownUser} key={unknownUser}>
                    <StyledAvatar alt={unknownUser} src={`${githubBaseURL}ghost.png`} />
                  </Tooltip>
                )}
              </StyledAvatarGroup>
            </StyledCommitRange>
          </Grid>
        </Grid>
      )}
    </StyledCommitContainer>
  );
};

export default CommitInformation;
