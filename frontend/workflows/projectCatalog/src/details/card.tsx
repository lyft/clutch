import React from "react";
import { Card, ClutchError, Error, Grid, styled, Typography } from "@clutch-sh/core";
import { LinearProgress } from "@material-ui/core";

export type CardType = "Dynamic" | "Metadata";

export interface DetailsCard {
  type: CardType;
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
  /** Number in ms to refresh the data from fetchDataFn */
  reloadInterval?: number;
  /** Boolean representing whether the component should reload via the fetchDataFn */
  autoReload?: boolean;
  /** Given promise which will be used to initially fetch data and optionally reload on intervals */
  fetchDataFn?: () => Promise<unknown>;
  /** Function called when fetchDataFn returns successfully, returning the data */
  onSuccess?: (data: any) => void;
  /** Function called when fetchDataFn returns unsuccessfully, returning an error */
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

const BaseCard = ({
  children,
  title,
  titleIcon,
  endAdornment,
  loading,
  error,
  ...props
}: CardProps) => {
  const [isLoading, setIsLoading] = React.useState<boolean>(false);
  const [cardError, setCardError] = React.useState<ClutchError | undefined>(undefined);

  const fetchData = () => {
    const { fetchDataFn, onSuccess, onError } = props;

    if (fetchDataFn) {
      setIsLoading(true);

      fetchDataFn()
        .then(data => {
          if (onSuccess) {
            onSuccess(data);
          }
          setCardError(undefined);
        })
        .catch(err => {
          if (onError) {
            onError(err);
          }
          setCardError(err);
        })
        .finally(() => setIsLoading(false));
    }
  };

  React.useEffect(() => {
    const { autoReload = false, reloadInterval = 30000 } = props;
    let interval;

    fetchData();

    if (autoReload) {
      interval = setInterval(fetchData, reloadInterval);
    }

    return () => (interval ? clearInterval(interval) : undefined);
  }, []);

  return (
    <StyledCard container direction="row">
      <Grid container item direction="row" alignItems="flex-start">
        <CardTitle title={title} titleIcon={titleIcon} endAdornment={endAdornment} />
      </Grid>
      <CardBody loading={loading || isLoading} error={error || cardError}>
        {children}
      </CardBody>
    </StyledCard>
  );
};

const DynamicCard = (props: BaseCardProps) => <BaseCard type="Dynamic" {...props} />;

const MetaCard = (props: BaseCardProps) => <BaseCard type="Metadata" {...props} />;

export { DynamicCard, MetaCard, StyledCard };
