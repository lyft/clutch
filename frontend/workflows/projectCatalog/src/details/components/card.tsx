import React from "react";
import type { TypographyProps } from "@clutch-sh/core";
import { Card, ClutchError, Error, Grid, styled, Typography } from "@clutch-sh/core";
import type { GridProps } from "@mui/material";
import { LinearProgress, Theme } from "@mui/material";

export enum CardType {
  DYNAMIC = "Dynamic",
  METADATA = "Metadata",
}

export interface CatalogDetailsCard {
  type: CardType;
}

interface CardTitleProps {
  title?: string | Element | React.ReactNode;
  titleVariant?: TypographyProps["variant"];
  titleAlign?: GridProps["alignItems"];
  titleIcon?: React.ReactNode;
  titleIconAlign?: GridProps["alignItems"];
  endAdornment?: React.ReactNode;
}

interface CardBodyProps {
  children?: React.ReactNode;
  /** Manual Control of loading state */
  loading?: boolean;
  /** Option to hide the loading indicator */
  loadingIndicator?: boolean;
  /** Manual control of error state */
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
  onSuccess?: (data: unknown) => void;
  /** Function called when fetchDataFn returns unsuccessfully, returning an error */
  onError?: (error: ClutchError | undefined) => void;
}

interface CardProps extends CatalogDetailsCard, BaseCardProps {}

const StyledCard = styled(Card)({
  width: "100%",
  height: "100%",
  padding: "16px",
});

const StyledGrid = styled(Grid)({
  height: "fit-content",
});

const StyledProgressContainer = styled("div")(({ theme }: { theme: Theme }) => ({
  marginBottom: "8px",
  marginTop: "-12px",
  height: "4px",
  ".MuiLinearProgress-root": {
    backgroundColor: theme.palette.primary[400],
  },
  ".MuiLinearProgress-bar": {
    backgroundColor: theme.palette.primary[600],
  },
}));

const StyledTitle = styled(Grid)({
  textTransform: "capitalize",
});

const CardTitle = ({
  title,
  titleVariant = "h4",
  titleAlign = "flex-start",
  titleIcon,
  titleIconAlign = "center",
  endAdornment,
}: CardTitleProps) => (
  <Grid
    container
    direction="row"
    alignItems={titleAlign}
    justifyContent="space-between"
    flexWrap="nowrap"
  >
    {title && (
      <Grid container item xs spacing={1} alignItems={titleIconAlign}>
        {titleIcon && <Grid item>{titleIcon}</Grid>}
        <StyledTitle item>
          <Typography variant={titleVariant}>{title}</Typography>
        </StyledTitle>
      </Grid>
    )}
    {endAdornment && (
      <Grid
        container
        item
        direction="row"
        xs
        spacing={1}
        alignItems="center"
        justifyContent="flex-end"
      >
        {endAdornment}
      </Grid>
    )}
  </Grid>
);

const CardBody = ({ loading, loadingIndicator = true, error, children }: CardBodyProps) => (
  <StyledGrid container direction="column" flexWrap="nowrap">
    {loadingIndicator && loading && (
      <Grid item>
        <StyledProgressContainer>
          <LinearProgress color="secondary" />
        </StyledProgressContainer>
      </Grid>
    )}
    <Grid item>
      <>{error ? <Error subject={error} /> : children}</>
    </Grid>
  </StyledGrid>
);

const BaseCard = ({ loading, error, ...props }: CardProps) => {
  const [cardLoading, setCardLoading] = React.useState<boolean>(false);
  const [cardError, setCardError] = React.useState<ClutchError | undefined>(undefined);

  const fetchData = () => {
    const { fetchDataFn, onSuccess, onError } = props;

    if (fetchDataFn) {
      setCardLoading(true);

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
        .finally(() => setCardLoading(false));
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
    <StyledCard>
      <Grid container direction="column" flexWrap="nowrap" spacing={2}>
        <Grid item>
          <CardTitle {...props} />
        </Grid>
        <Grid item>
          <CardBody loading={loading || cardLoading} error={error || cardError} {...props} />
        </Grid>
      </Grid>
    </StyledCard>
  );
};

const DynamicCard = (props: BaseCardProps) => <BaseCard type={CardType.DYNAMIC} {...props} />;

const MetaCard = (props: BaseCardProps) => <BaseCard type={CardType.METADATA} {...props} />;

export type DetailCard = CatalogDetailsCard | typeof DynamicCard | typeof MetaCard;

export { DynamicCard, MetaCard };
