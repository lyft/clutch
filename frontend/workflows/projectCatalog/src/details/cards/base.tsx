import React from "react";
import { Card, ClutchError, Error, Grid, Link, styled, Typography } from "@clutch-sh/core";
import { faClock } from "@fortawesome/free-regular-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { LinearProgress } from "@material-ui/core";
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

export interface BaseProjectCardProps {
  text: string;
  icon?: React.ReactNode;
  endAdornment?: React.ReactNode;
  loadData?: () => Promise<any>;
  setData?: (data: any) => void;
  loading?: boolean;
  error?: ClutchError | undefined;
  reloadInterval?: number;
  autoReload?: boolean;
}

export interface ExtendedProjectCardProps extends BaseProjectCardProps {
  children?: React.ReactNode;
}

interface BaseCardProps {
  interval?: number;
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

const TitleRow = ({ text, icon, endAdornment }: BaseProjectCardProps) => (
  <>
    {icon && (
      <Grid item xs={1}>
        {icon}
      </Grid>
    )}
    <Grid item xs={8}>
      <Typography variant="h4">{text}</Typography>
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

export const BaseCard = ({
  children,
  text,
  icon,
  endAdornment,
  loading,
  error,
}: ExtendedProjectCardProps) => {
  return (
    <StyledCard container direction="row">
      <Grid container item direction="row" alignItems="flex-start">
        <TitleRow text={text} icon={icon} endAdornment={endAdornment} />
      </Grid>
      <StyledRow>
        {loading && (
          <StyledProgressContainer>
            {loading && <LinearProgress color="secondary" />}
          </StyledProgressContainer>
        )}
      </StyledRow>
      {/* {children} */}
      {error ? <Error subject={error} /> : children}
    </StyledCard>
  );
};

class BaseCardComponent extends React.Component<ExtendedProjectCardProps, BaseCardState> {
  static displayName = "BaseCard";

  constructor(props: ExtendedProjectCardProps) {
    super(props);
    const { loadData, autoReload, loading = false, error, reloadInterval = 30000 } = this.props;
    this.state = {
      error,
      loading,
      reloadInterval,
      interval: undefined,
      data: undefined,
    };
    console.log(this.props);

    if (loadData && autoReload) {
      this.fetchData(loadData);
      this.setPromiseInterval();
    }
  }

  componentDidCatch(baseError) {
    this.setState(state => ({ ...state, baseError }));
  }

  componentWillUnmount() {
    if (this.state.interval) {
      clearInterval(this.state.interval);
    }
  }

  fetchData(promise: () => Promise<any>) {
    this.setState(state => ({ ...state, loading: true }));
    console.log("Calling fetch data");

    promise()
      .then(data => {
        if (this.props.setData) {
          this.props.setData(data);
        }
      })
      .catch(error => this.setState(state => ({ ...state, error })))
      .finally(() => this.setState(state => ({ ...state, loading: false })));
  }

  setPromiseInterval() {
    const { loadData, autoReload } = this.props;
    if (loadData && autoReload) {
      if (this.state.interval) {
        clearInterval(this.state.interval);
        this.setState(state => ({ ...state, interval: undefined }));
      }

      const interval = setInterval(() => this.fetchData(loadData), this.state.reloadInterval);
      this.setState(state => ({ ...state, interval }));
    }
  }

  render() {
    return <></>;
  }
}

export default BaseCardComponent;
