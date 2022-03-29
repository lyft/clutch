import React from "react";
import { Card, ClutchError, Error, Grid, styled, Typography } from "@clutch-sh/core";
import { LinearProgress } from "@material-ui/core";

export type CardTypes = "Dynamic" | "Meta";

export interface DetailsCard {
  type: CardTypes;
}

interface CardTitleProps {
  title: string;
  titleIcon?: React.ReactNode;
  endAdornment?: React.ReactNode;
}

interface CardBodyProps {
  children?: React.ReactNode;
  loading?: boolean;
  error?: ClutchError;
}

interface BaseCardProps extends CardTitleProps, CardBodyProps {
  reloadInterval?: number;
  autoReload?: boolean;
  fetchDataFn?: () => Promise<any>;
  onSuccess?: (data: any) => void;
  onError?: (error: ClutchError | undefined) => void;
}

interface CardProps extends DetailsCard, BaseCardProps {}

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

const StyledRow = styled(Grid)({
  marginBottom: "15px",
});

const CardTitle = ({ title, titleIcon, endAdornment }: CardTitleProps) => (
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

const CardBody = ({ loading, error, children }: CardBodyProps) => (
  <>
    <StyledRow>
      {loading && (
        <StyledProgressContainer>
          <LinearProgress color="secondary" />
        </StyledProgressContainer>
      )}
    </StyledRow>
    {error ? <Error subject={error} /> : children}
  </>
);

const BaseCard = (props: CardProps) => {
  const [isLoading, setIsLoading] = React.useState<boolean>(false);
  const [clutchError, setClutchError] = React.useState<ClutchError | undefined>(undefined);
  const { children, title, titleIcon, endAdornment, loading, error } = props;

  const fetchData = (promise: () => Promise<any>) => {
    const { onSuccess, onError } = props;

    setIsLoading(true);

    promise()
      .then(data => {
        if (onSuccess) {
          onSuccess(data);
        }
      })
      .catch(err => {
        if (onError) {
          onError(err);
        }
        setClutchError(err);
      })
      .finally(() => setIsLoading(false));
  };

  React.useEffect(() => {
    const { fetchDataFn, autoReload, reloadInterval = 30000 } = props;
    let interval;

    if (fetchDataFn) {
      fetchData(fetchDataFn);

      if (autoReload) {
        interval = setInterval(() => fetchData(fetchDataFn), reloadInterval);
      }
    }

    return () => (interval ? clearInterval(interval) : undefined);
  }, []);

  return (
    <StyledCard container direction="row">
      <Grid container item direction="row" alignItems="flex-start">
        <CardTitle title={title} titleIcon={titleIcon} endAdornment={endAdornment} />
      </Grid>
      <CardBody loading={loading || isLoading} error={error || clutchError}>
        {children}
      </CardBody>
    </StyledCard>
  );
};

const DynamicCard = (props: BaseCardProps) => <BaseCard type="Dynamic" {...props} />;

const MetaCard = (props: BaseCardProps) => <BaseCard type="Meta" {...props} />;

export { DynamicCard, MetaCard, StyledCard };
