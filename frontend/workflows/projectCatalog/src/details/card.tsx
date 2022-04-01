import React from "react";
import { Card, ClutchError, Error, Grid, styled, Typography } from "@clutch-sh/core";
import { LinearProgress } from "@material-ui/core";

enum CardType {
  DYNAMIC = "Dynamic",
  METADATA = "Metadata",
}

export interface CatalogDetailsCard {
  type: CardType;
}

interface CardTitleProps {
  title?: string;
  titleIcon?: React.ReactNode;
  endAdornment?: React.ReactNode;
}

interface CardBodyProps {
  children?: React.ReactNode;
  /** Manual Control of loading state */
  loading?: boolean;
  /** Manual control of error state */
  error?: ClutchError;
}

interface BaseCardProps extends CardTitleProps, CardBodyProps {
  /** Number in ms to refresh the data from fetchDataFn */
  reloadInterval?: number;
  /** Boolean representing whether the component should reload via the fetchDataFn */
  autoReload?: boolean;
  /** Optionally disable loading indicator */
  loadingIndicator?: boolean;
  /** Given promise which will be used to initially fetch data and optionally reload on intervals */
  fetchDataFn?: () => Promise<unknown>;
  /** Function called when fetchDataFn returns successfully, returning the data */
  onSuccess?: (data: unknown) => void;
  /** Function called when fetchDataFn returns unsuccessfully, returning an error */
  onError?: (error: ClutchError | undefined) => void;
}

interface CardProps extends CatalogDetailsCard, BaseCardProps {}

const StyledCard = styled(Card)({
  width: "100%",
  height: "fit-content",
  padding: "16px",
});

const StyledProgressContainer = styled("div")({
  marginBottom: "8px",
  marginTop: "-12px",
  height: "4px",
  ".MuiLinearProgress-root": {
    backgroundColor: "rgb(194, 200, 242)",
  },
  ".MuiLinearProgress-bar": {
    backgroundColor: "#3548D4",
  },
});

const StyledTitle = styled(Grid)({
  textTransform: "capitalize",
});

const StyledTitleContainer = styled(Grid)({
  marginBottom: "8px",
});

const BodyContainer = styled("div")({
  paddingLeft: "4px",
});

const CardTitle = ({ title, titleIcon, endAdornment }: CardTitleProps) => (
  <>
    {title && (
      <StyledTitleContainer container item xs={endAdornment ? 9 : 12} spacing={1}>
        {titleIcon && <Grid item>{titleIcon}</Grid>}
        <StyledTitle item>
          <Typography variant="h4">{title}</Typography>
        </StyledTitle>
      </StyledTitleContainer>
    )}
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
    {loading && (
      <StyledProgressContainer>
        <LinearProgress color="secondary" />
      </StyledProgressContainer>
    )}
    <BodyContainer>{error ? <Error subject={error} /> : children}</BodyContainer>
  </>
);

const BaseCard = ({
  children,
  title,
  titleIcon,
  endAdornment,
  loading,
  loadingIndicator = true,
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
      <CardBody loading={loadingIndicator && (loading || isLoading)} error={error || cardError}>
        {children}
      </CardBody>
    </StyledCard>
  );
};

const DynamicCard = (props: BaseCardProps) => <BaseCard type={CardType.DYNAMIC} {...props} />;

const MetaCard = (props: BaseCardProps) => <BaseCard type={CardType.METADATA} {...props} />;

export { CardType, DynamicCard, MetaCard };
