export const tickFormatterFunc = timeStamp => {
  const date = new Date(timeStamp);
  return date.toLocaleTimeString();
};

// Edge ratio refers to the amount that the min will be subtracted to, and the
// amount the max will be added to.
export const calculateDomainEdges = (data, dataKey: string, edgeRatio: number) => {
  // Get the max and min of the domain, then calculate percent`out from each edge.
  const min = Math.min(...data.map(d => d[dataKey]));
  const max = Math.max(...data.map(d => d[dataKey]));
  if (max === min) {
    const minEdge = min - min * edgeRatio;
    const maxEdge = max + max * edgeRatio;
    return [minEdge, maxEdge];
  }
  const minEdge = min - (max - min) * edgeRatio;
  const maxEdge = max + (max - min) * edgeRatio;
  return [minEdge, maxEdge];
};
