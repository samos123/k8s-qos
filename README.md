Dynamic Bandwidth QoS in K8s
============================
[![Build Status](https://travis-ci.org/samos123/k8s-qos.svg?branch=master)](https://travis-ci.org/samos123/k8s-qos)

Multi-tenant K8s clusters often see the issue of one pod taking up all
the available bandwidth on a node. This might cause loss of service for
other pods that run on the same node. This is considered a noisy neighbor
that is taking up all the available bandwidth.
This project introduces a new dynamic bandwidth QoS mechanism for K8s to
solve the bandwidth heavy noisy neighbor problem by being able to identify
high bandwidth/network throughput pods and automatically applying limits or
alerting operators.

Main features will include:
*  Analyze network throughput of pods on a node
*  Sent an alert when a pod uses most of the bandwidth
*  Automatically apply bandwidth limits to pods that regularly use most of the bandwidth
*  Support different automatic QoS strategies


Bandwidth QoS strategies
------------------------
Bandwidth isn't a resource that you nessecarily want to limit unless you
absolutely have to. Hence you don't want to treat it like CPU and/or memory
resources, where you generally don't oversubscribe.

For example let's take the example of running 10 pods on 1 node, would you
want to limit each pod 10% of the available bandwidth on a node? Probably not,
instead you might want to prevent that a single pod can never take more than
70% of the available bandwidth on a node. However, what if you run only 1 pod
on a node then you wouldn't want any bandwidth limits.

So as part of this project, the following strategies will be available:
*  Pod amount limiter (Default)
*  Actual network throughput based limiter

### Pod Amount limiter
The pod amount limiter looks at the current amount of pods and automatically
sets bandwidth QoS limits. The table below explains what actions the
limter takes. The #Pods column is the amount of pods on a node. The limit
column is limit of total available bandwidth of the node.

| #Pods | Limit    |
|-------|----------|
| 1     | No limit |
| 2-3   | 80%      |
| 4-5   | 70%      |
| 6-10  | 60%      |
| 10+   | 40%      |


### Actual network throughput based limiter
This limiter will look at current network throughput of nodes and
automatically apply limits to pods that regularly utilize 80% or
more of the bandwidth.

Why another project instead of the CNI bandwidth plugin?
--------------------------------------------------------
Today K8s doesn't treat bandwidth as a resource that can be limited or used
as a resource for scheduling pods. There is a [bandwidth plugin](https://kubernetes.io/docs/concepts/extend-kubernetes/compute-storage-net/network-plugins/#support-traffic-shaping)
that allows you set bandwidth ingress and egress limits to pods, however it has
the following limitations:
*  Requires the operator to know what bandwidth limits should be set
*  Installation of a CNI plugin, which may not be possible
*  Unable to automatically adjust based on traffic patterns

