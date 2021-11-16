package xds

import (
	"fmt"

	gcpType "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	serverexperimentation "github.com/lyft/clutch/backend/api/chaos/serverexperimentation/v1"
)

const (
	UpstreamEnforcing   = "upstreamEnforcing"
	DownstreamEnforcing = "downstreamEnforcing"
	UnknownEnforcing    = "unknown"
)

func GetClusterPair(httpFaultConfig *serverexperimentation.HTTPFaultConfig) (string, string, error) {
	var downstream, upstream string

	switch httpFaultConfig.GetFaultTargeting().GetEnforcer().(type) {
	case *serverexperimentation.FaultTargeting_DownstreamEnforcing:
		downstreamEnforcing := httpFaultConfig.GetFaultTargeting().GetDownstreamEnforcing()

		switch downstreamEnforcing.GetDownstreamType().(type) {
		case *serverexperimentation.DownstreamEnforcing_DownstreamCluster:
			downstream = downstreamEnforcing.GetDownstreamCluster().GetName()
		default:
			return "", "", fmt.Errorf("unknown downstream type of downstream enforcing %v", downstreamEnforcing.GetDownstreamType())
		}

		switch downstreamEnforcing.GetUpstreamType().(type) {
		case *serverexperimentation.DownstreamEnforcing_UpstreamCluster:
			upstream = downstreamEnforcing.GetUpstreamCluster().GetName()
		default:
			return "", "", fmt.Errorf("unknown upstream type of downstream enforcing %v", downstreamEnforcing.GetUpstreamType())
		}

	case *serverexperimentation.FaultTargeting_UpstreamEnforcing:
		upstreamEnforcing := httpFaultConfig.GetFaultTargeting().GetUpstreamEnforcing()

		switch upstreamEnforcing.GetDownstreamType().(type) {
		case *serverexperimentation.UpstreamEnforcing_DownstreamCluster:
			downstream = upstreamEnforcing.GetDownstreamCluster().GetName()
		default:
			return "", "", fmt.Errorf("unknown downstream type of upstream enforcing %v", upstreamEnforcing.GetDownstreamType())
		}

		switch upstreamEnforcing.GetUpstreamType().(type) {
		case *serverexperimentation.UpstreamEnforcing_UpstreamCluster:
			upstream = upstreamEnforcing.GetUpstreamCluster().GetName()
		default:
			return "", "", fmt.Errorf("unknown upstream type of upstream enforcing %v", upstreamEnforcing.GetUpstreamType())
		}

	default:
		return "", "", fmt.Errorf("unknown enforcer %v", httpFaultConfig.GetFaultTargeting())
	}

	return upstream, downstream, nil
}

func GetEnforcer(httpFaultConfig *serverexperimentation.HTTPFaultConfig) string {
	switch httpFaultConfig.GetFaultTargeting().GetEnforcer().(type) {
	case *serverexperimentation.FaultTargeting_DownstreamEnforcing:
		return DownstreamEnforcing
	case *serverexperimentation.FaultTargeting_UpstreamEnforcing:
		return UpstreamEnforcing
	default:
		return UnknownEnforcing
	}
}

func GetEnforcingCluster(httpFaultConfig *serverexperimentation.HTTPFaultConfig) (string, error) {
	upstream, downstream, err := GetClusterPair(httpFaultConfig)
	if err != nil {
		return "", err
	}

	enforcer := GetEnforcer(httpFaultConfig)
	switch enforcer {
	case DownstreamEnforcing:
		return downstream, nil
	case UpstreamEnforcing:
		return upstream, nil
	default:
		return "", fmt.Errorf("unknown enforcing type %s", enforcer)
	}
}

// GetHTTPFaultPercentage is a helper to translate the percentage field of a
// HTTPFaultConfig into a FractionalPercent. See:
// https://www.envoyproxy.io/docs/envoy/latest/api-v3/type/v3/percent.proto#type-v3-fractionalpercent
func GetHTTPFaultPercentage(httpFaultConfig *serverexperimentation.HTTPFaultConfig) (*gcpType.FractionalPercent, error) {
	var expPercentage *serverexperimentation.FaultPercentage
	switch httpFaultConfig.GetFault().(type) {
	case *serverexperimentation.HTTPFaultConfig_AbortFault:
		expPercentage = httpFaultConfig.GetAbortFault().GetPercentage()
	case *serverexperimentation.HTTPFaultConfig_LatencyFault:
		expPercentage = httpFaultConfig.GetLatencyFault().GetPercentage()
	default:
		return nil, fmt.Errorf("unknown enforcer %v", httpFaultConfig)
	}

	var denom gcpType.FractionalPercent_DenominatorType
	switch expPercentage.GetDenominator() {
	case serverexperimentation.FaultPercentage_DENOMINATOR_TEN_THOUSAND:
		denom = gcpType.FractionalPercent_TEN_THOUSAND
	case serverexperimentation.FaultPercentage_DENOMINATOR_MILLION:
		denom = gcpType.FractionalPercent_MILLION
	default:
		denom = gcpType.FractionalPercent_HUNDRED
	}

	return &gcpType.FractionalPercent{
		Numerator:   expPercentage.GetPercentage(),
		Denominator: denom,
	}, nil
}
