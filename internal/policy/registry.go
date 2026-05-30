package policy

import "fmt"

// PlacementPolicyFactory is a function that builds a PlacementPolicy from opts
type PlacementPolicyFactory func(opts PolicyOpts) (PlacementPolicy, error)

var placementPolicyFactories = map[string]PlacementPolicyFactory{}

// Register adds a named policy factor. call this in each policy's init()
func Register(name string, f PlacementPolicyFactory) {
	placementPolicyFactories[name] = f
}

func Build(name string, opts PolicyOpts) (PlacementPolicy, error) {
	f, ok := placementPolicyFactories[name]
	if !ok {
		return nil, fmt.Errorf("[policy.Build]: unknown policy %q", name)
	}
	return f(opts)
}

func Available() []string {
	names := make([]string, 0, len(placementPolicyFactories))
	for name := range placementPolicyFactories {
		names = append(names, name)
	}
	return names
}
