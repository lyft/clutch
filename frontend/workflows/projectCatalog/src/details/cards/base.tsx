import React from "react";
import { Card, ClutchError, Error, Grid, Link, styled, Typography } from "@clutch-sh/core";
import { faClock } from "@fortawesome/free-regular-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { LinearProgress } from "@material-ui/core";

import type { DetailsCardTypes } from "../..";
import { EventTime, setMilliseconds } from "../helpers";

const StyledCard = styled(Card)({
  width: "100%",
  height: "fit-content",
  padding: "15px",
});

const StyledProgressContainer = styled("div")({
  height: "4px",
  ".MuiLinearProgress-root": {
    backgroundColor: "rgb(194, 200, 242)",
  },
  ".MuiLinearProgress-bar": {
    backgroundColor: "#3548D4",
  },
});

const StyledLink = styled(Link)({
  whiteSpace: "nowrap",
});

const StyledRow = styled(Grid)({
  marginBottom: "15px",
});

const LinkText = ({ text, link }: { text: string; link?: string }) => {
  const returnText = <Typography variant="body2">{text}</Typography>;

  if (link && text) {
    return <StyledLink href={link}>{returnText}</StyledLink>;
  }

  return returnText;
};

interface TitleRowProps {
  title: string;
  titleIcon?: React.ReactNode;
  endAdornment?: React.ReactNode;
}

export interface BaseCardProps extends TitleRowProps {
  type?: DetailsCardTypes;
  loading?: boolean;
  error?: ClutchError | undefined;
  reloadInterval?: number;
  autoReload?: boolean;
  fetchDataFn?: () => Promise<any>;
  onSuccess?: (data: any) => void;
  onError?: (error: ClutchError | undefined) => void;
}

export interface ProjectCardProps extends BaseCardProps {
  children?: React.ReactNode;
}

interface BaseCardState {
  data?: any;
  interval?: number;
  loading?: boolean;
  reloadInterval?: number;
  error?: ClutchError | undefined;
}

const LastEvent = ({ time }: { time: number }) => (
  <>
    {time && (
      <>
        <Grid item>
          <FontAwesomeIcon icon={faClock} />
        </Grid>
        <Grid item>
          <Typography variant="body4">
            <EventTime date={setMilliseconds(time)} /> ago
          </Typography>
        </Grid>
      </>
    )}
  </>
);

const TitleRow = ({ title, titleIcon, endAdornment }: TitleRowProps) => (
  <>
    {titleIcon && (
      <Grid item xs={1}>
        {titleIcon}
      </Grid>
    )}
    <Grid item xs={8}>
      <Typography variant="h4">{title}</Typography>
    </Grid>
    {endAdornment && (
      <Grid
        container
        item
        direction="row"
        xs={3}
        spacing={1}
        alignItems="center"
        justify="flex-end"
      >
        {endAdornment}
      </Grid>
    )}
  </>
);

const BaseCard = ({
  children,
  title,
  titleIcon,
  endAdornment,
  loading,
  error,
}: ProjectCardProps) => {
  return (
    <StyledCard container direction="row">
      <Grid container item direction="row" alignItems="flex-start">
        <TitleRow title={title} titleIcon={titleIcon} endAdornment={endAdornment} />
      </Grid>
      <StyledRow>
        {loading && (
          <StyledProgressContainer>
            {loading && <LinearProgress color="secondary" />}
          </StyledProgressContainer>
        )}
      </StyledRow>
      {error ? <Error subject={error} /> : children}
    </StyledCard>
  );
};

class BaseCardComponent extends React.Component<BaseCardProps, BaseCardState> {
  constructor(props: BaseCardProps) {
    super(props);
    const { loading = false, error, reloadInterval = 30000 } = this.props;
    this.state = {
      error,
      loading,
      reloadInterval,
      interval: undefined,
      data: undefined,
    };
  }

  componentDidMount() {
    const { fetchDataFn, autoReload } = this.props;

    if (fetchDataFn) {
      this.fetchData(fetchDataFn);

      if (autoReload) {
        this.setPromiseInterval();
      }
    }
  }

  componentDidCatch(error) {
    this.setState(state => ({ ...state, error }));
  }

  componentWillUnmount() {
    const { interval } = this.state;

    if (interval) {
      clearInterval(interval);
    }
  }

  setPromiseInterval() {
    const { fetchDataFn, autoReload } = this.props;
    const { interval, reloadInterval } = this.state;
    if (fetchDataFn && autoReload) {
      if (interval) {
        clearInterval(interval);
        this.setState(state => ({ ...state, interval: undefined }));
      }

      const newInterval = setInterval(() => this.fetchData(fetchDataFn), reloadInterval);
      this.setState(state => ({ ...state, interval: newInterval }));
    }
  }

  fetchData(promise: () => Promise<any>) {
    const { onSuccess, onError } = this.props;
    this.setState(state => ({ ...state, loading: true }));

    promise()
      .then(data => {
        if (onSuccess) {
          onSuccess(data);
        }
      })
      .catch(error => {
        if (onError) {
          onError(error);
        }
        this.setState(state => ({ ...state, error }));
      })
      .finally(() => this.setState(state => ({ ...state, loading: false })));
  }

  render() {
    return <></>;
  }
}

export { BaseCard, LastEvent, LinkText, StyledCard, StyledLink, StyledRow };

export default BaseCardComponent;
