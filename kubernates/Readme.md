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
    
🔹 Node
    A Node is a physical or virtual machine in the cluster.
    Each node contains:
    - Kubelet (agent)
    - Container Runtime (e.g., Docker, containerd)
    - Kube Proxy (networking)
    
🔹 Pod
    The smallest deployable unit in Kubernetes.
    A Pod can contain one or more tightly coupled containers that share the same network and storage.
    
🔹 Deployment
  Defines the desired state for Pods and manages their lifecycle (e.g., rolling updates, rollbacks).
  Ensures the specified number of Pods are running at any given time.
  
🔹 Service
  Exposes an application running on a set of Pods as a network service.
  
  Types include:
  - ClusterIP: Internal communication
  - NodePort: Exposes service on each Node’s IP
  - LoadBalancer: External access via cloud provider load balancer

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
![image](https://github.com/user-attachments/assets/9e8b19c6-7fe1-40c9-b499-c5b7d4103544)

### Control Plane
**Component of Control Plane**
- Kube-apiserver
- etcd
- kube-scheduler
- controller-manager

1. Kube-apiserver(for all commuinication)
   - it ineracts direclty with user (using manifest file of .yml or .json)
   - this kube-apiserver is scale automatically as per load
   - it is frontend of control plane

2. etcd(database)
   - stores metadata and status of cluster
   - it it consistant and high-avaiable store(key-value)
     Features :
     1. Fully replicates: the entire state is available on everynode in the cluster
     2. secure: implements automatics tls with optimal client-certificate authentication
     3. fast: benchmarked at 10,000 writes per second

3. kube-scheduler(action)
   - when user make request for the creation & management of pods, kube-scheduler is going to take action on these requests
   - handles pod creation and management
   - kube-scheduler match/assign any node to create and run pods
   - scheduler watches for newly created pods that have no node assigned. for every pod that the scheduler discovers, the scheduler becomes responsible for finding best node for that pod to run on it.
   - Scheduler get the information for hardware configuration from config/manifest files and schedules the pods on nodes accordingly
  
4. controller-manager
   - make sure actual state of cluster mateches to desired state
   - two possible cloise for controll manage:
     1. cloud-controll-manager(work on aws/gcp/azure)
     2. kube-controll-manager(work on non cloud/on-premisis environments)
   
   **Components on master that runs controller**
   1. node-controller: for checking the cloud provider to determine if a node has been detected in the cloud after it stops resonding
   2. route-controller: responsible for setting up network routes in cloud
   3. service-controller: responsible for load balancers on your cloud against services of type load balancer
   4. volume-controller: for creating, attaching and mouting volumes and interacting with the cloud providers to orchestrate volume
      
### Nodes  

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
