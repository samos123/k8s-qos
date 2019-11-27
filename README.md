Dynamic Bandwidth QoS in K8s
============================

Multi-tenant K8s clusters often see the issue of one pod taking up all
the available bandwidth on a node. This might cause loss of service for
other pods that run on the same node. This is considered a noisy neighbor
that is taking up all the available bandwidth.
This project introduces a new dynamic bandwidth QoS mechanism for K8s to
solve the bandwidth heavy noisy neighbor problem by being able to identify
high bandwidth pods and automatically applying limits or alerting operators.

Main features will include:
*  Analyze throughput of pods on a node
*  Sent an alert when a pod uses most of the bandwidth
*  Automatically apply bandwidth limits to pods that regularly use most of the bandwidth

Why another project instead of the CNI bandwidth plugin?
--------------------------------------------------------
Today K8s doesn't treat bandwidth as a resource that can be limited or used
as a resource for scheduling pods. There is a [bandwidth plugin](https://kubernetes.io/docs/concepts/extend-kubernetes/compute-storage-net/network-plugins/#support-traffic-shaping)
that allows you set bandwidth ingress and egress limits to pods, however it has
the following limitations:
*  Requires the operator to know what bandwidth limits should be set
*  Installation of a CNI plugin, which may not be possible
*  Unable to automatically adjust based on traffic patterns

