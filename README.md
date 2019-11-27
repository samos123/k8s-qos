Dynamic Bandwidth QoS in K8s
============================

Multi-tenant K8s clusters often see the issue of one pod taking up all
the available bandwidth on a node. This might cause loss of service for
other pods that run on the same node. This is considered a noisy neighbor
that is taking up all the available bandwidth.

Today K8s doesn't treat bandwidth as a resource that can be limited or used
as a resource for scheduling pods. There is a bandwidth plugin that
allows you set bandwidth ingress and egress limits to pods, however it has
the following limitations:
*  Requires the operator to know what bandwidth limits should be set
*  Installation of a CNI plugin that may not work on all K8s clusters

This project introduces a new dynamic bandwidth QoS mechanism for K8s.
