# cnsim
for simulating container networking

## spec

- scenarios:
  - steady state, no churn
  - when a new app is pushed
  - when a host goes down
  - when rolling out a OS upgrade
  - when an entire AZ goes down

- inputs:
  - numHosts
  - numApps
  - distAppSize
  - probReflexive # probably that an app connects to itself
  - distAppDegree # distribution of degrees on graph of apps (# other apps it connects to)

- outputs:
  - medianAppsPerHost
  - distHostDegree # distribution of degrees on graph of hosts
  - distHostRouteTableSize

- for the future:
  - multi-tenant
  - non-uniform app sizes, with packing
  - more realistic simulation of the auction
  - options for affinity / anti-affinity
