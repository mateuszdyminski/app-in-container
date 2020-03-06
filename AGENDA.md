# Agenda

* Building Container
  * Choosing best Docker base image
  * Multistage builds
  * Layering as to improve build time
  * Google distroless
* Security
  * Nonroot
  * The less soft image contains the better security
  * Scanning vulnerabilities
  * Static containers analysis
* Intelligent Health Checks - useful for K8s
  * Liveness probes
  * Readiness probes
  * Startup probes
  * Exmaple + demo on K8s
* Graceful Shutdown – handling SIGINT and other signals
  * Example in go how to handle termination
  * Customer Docker signal
  * Demo with graceful + rolling update ?
* Exposing metrics
  * Prometheus integration + demo
  * Java + sidecar
* Application configuration – reloading app configuration without downtime
  * Reload pattern with sidecar
