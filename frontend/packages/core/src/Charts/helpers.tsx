export const localTimeFormatter = timeStamp => {
  const date = new Date(timeStamp);
  return date.toLocaleTimeString();
};

const getMinAndMaxOfRangeUsingKey = (data, key) => {
  const min = Math.min(...data.filter(d => d[key]).map(d => d[key]));
  const max = Math.max(...data.filter(d => d[key]).map(d => d[key]));
  return { min, max };
};

// Edge ratio refers to the multiplicative part of the amount that will be added to the max,
// and subtracted from the min.
export const calculateDomainEdges = (data, dataKey: string, edgeRatio: number) => {
  // Get the max and min of the domain, then calculate a certain amount`out from each edge.
  const { min, max } = getMinAndMaxOfRangeUsingKey(data, dataKey);
  if (edgeRatio <= 0) {
    return [min, max];
  }
  if (max === min) {
    const minEdge = min - min * edgeRatio;
    const maxEdge = max + max * edgeRatio;
    return [minEdge, maxEdge];
  }
  const edgeVal = (max - min) * edgeRatio;
  const minEdge = min - edgeVal;
  const maxEdge = max + edgeVal;
  return [minEdge, maxEdge];
};

export const calculateTicks = (data, dataKey: string) => {};
