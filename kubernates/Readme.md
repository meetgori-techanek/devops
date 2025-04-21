# Kubernates

## History
- Kubernetes is an open-source container orchestration tool. 
- Originally developed by Google in 2014 for internal container management, where it was initially called Borg (not Coorg). 
- Written in Go (Golang).
- Now maintained by the Cloud Native Computing Foundation (CNCF).

## Why Use
- Kubernetes is widely used in modern DevOps practices because it simplifies the deployment, scaling, and management of containerized applications. Key reasons to use Kubernetes include:
- Automation: Automates deployment, scaling, and operations of application containers.
- Scalability: Easily scale applications up or down based on demand, either manually or automatically.
- High Availability: Ensures application uptime through self-healing and fault tolerance.
- Portability: Works across different cloud providers (AWS, GCP, Azure) and on-premises environments.
- Consistency: Provides a consistent development, testing, and production environment.
- DevOps Friendly: Integrates well with CI/CD pipelines, promoting agile and DevOps workflows.
- Resource Optimization: Efficiently uses infrastructure resources via scheduling and resource allocation.
  
## Pros and Cons
- **Pros:**
  - Orchestration: Manages clusters of nodes and containers, even across multiple clouds.
  - Autoscaling: Supports both horizontal (add more pods) and vertical (increase pod resources) scaling.
  -  Load Balancing: Distributes traffic across services to maintain performance.
  - Platform Independent: Runs on any infrastructure (cloud or on-prem).
  - Fault Tolerance: Automatically restarts failed pods, replaces them, and handles node failures.
  - Rollback: Revert to previous versions of applications if something goes wrong.
  - Health Monitoring: Actively checks the health of applications and replaces unhealthy containers.
  - Batch Execution: Handles batch jobs that can be run once, sequentially, or in parallel.
- **Cons:**
  - Complexity: Steep learning curve for beginners.
  - Resource Intensive: Requires significant system resources and configuration.
  - Overhead: Not ideal for small-scale applications or simple deployments.
  - Security: Needs careful setup and management to avoid security vulnerabilities.

## Fundamentals
🔹 Cluster
A Kubernetes Cluster consists of:

Master Node (Control Plane): Manages the cluster and handles scheduling, scaling, and deployment.

Worker Nodes: Where the actual applications (containers) run.

🔹 Pod
The smallest deployable unit in Kubernetes.

A Pod can contain one or more tightly coupled containers that share the same network and storage.

🔹 Node
A Node is a physical or virtual machine in the cluster.

Each node contains:

Kubelet (agent)

Container Runtime (e.g., Docker, containerd)

Kube Proxy (networking)

🔹 Deployment
Defines the desired state for Pods and manages their lifecycle (e.g., rolling updates, rollbacks).

Ensures the specified number of Pods are running at any given time.

🔹 Service
Exposes an application running on a set of Pods as a network service.

Types include:

ClusterIP: Internal communication

NodePort: Exposes service on each Node’s IP

LoadBalancer: External access via cloud provider load balancer

🔹 Namespace
Logical partitions within a cluster for isolating resources and managing permissions.

🔹 ConfigMap & Secret
ConfigMap: External configuration (non-sensitive).

Secret: Stores sensitive information like passwords or tokens securely.

🔹 Volume
Persistent storage attached to a Pod.

Can be cloud-based (e.g., AWS EBS) or local.

🔹 Ingress
Manages external access to services, typically via HTTP.

Supports load balancing, SSL termination, and name-based virtual hosting.

## Architecture
- Overview of the system architecture and components.

## Installation
- Step-by-step instructions for installing the tool.

## Commands
- Common commands and their usage.

## Best Practices
- Recommended practices for effective use.

## Security Checks
- Guidelines for ensuring security while using the tool.

## Secret Management
- Strategies for managing sensitive information securely.
