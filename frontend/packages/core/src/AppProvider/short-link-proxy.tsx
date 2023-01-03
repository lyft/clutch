import React from "react";
import { useParams } from "react-router-dom";
import type { clutch as IClutch } from "@clutch-sh/api";

import type { HydrateState } from "../Contexts/workflow-storage-context/types";
import { useNavigate } from "../Navigation/navigation";
import { client } from "../Network";
import type { ClutchError } from "../Network/errors";

/**
 * The base for a short link route
 */
export const ShortLinkBaseRoute = "goto";

/**
 * Will return a ShortLink route
 * @param origin the windows origin
 * @param hash the hash to use for the short link
 * @returns string
 */
export const generateShortLinkRoute = (origin: string, hash: string) =>
  `${origin}/${ShortLinkBaseRoute}/${hash}`;

interface ShortLinkProps {
  hydrate: (data: HydrateState) => void;
  onError: (error: ClutchError) => void;
  setLoading: (loading: boolean) => void;
}

const fetchData = async (hash, hydrate, onError) => {
  const requestData: IClutch.shortlink.v1.IGetRequest = { hash };

  return client
    .post("/v1/shortlink/get", requestData)
    .then(response => {
      const { path, state } = response.data as IClutch.shortlink.v1.IGetResponse;

      hydrate({ hash, state });
      return path;
    })
    .catch((error: ClutchError) => {
      onError(error);
      return "/";
    });
};

/**
 * Component that will be present for a route which will look for a short link hash
 * If found
 * - It will set a loading state
 * - Then it will call down to the API with the hash and ask for any data pertaining to it
 * - If the API call is successful
 *   - It will use the given hydrate function to send the returned state off to the StorageContext
 *   - It will navigate to the route given in the returned state
 * - If the API call is not successful
 *   - It will leave a warning message in the console
 *   - Then navigate back to the home page
 * - Then it will remove the loading state
 */
const ShortLinkProxy = ({ hydrate, onError, setLoading }: ShortLinkProps) => {
  const { hash } = useParams();
  const navigate = useNavigate();

  React.useEffect(() => {
    if (hash) {
      setLoading(true);
      (async function loadSL() {
        const path = await fetchData(hash, hydrate, onError);

        navigate(path);
        setLoading(false);
      })();
    }
  }, [hash]);

  // currently return null so that nothing is rendered
  return null;
};

export default ShortLinkProxy;
