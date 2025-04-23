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
![image](https://kubernetes.io/images/docs/kubernetes-cluster-architecture.svg)


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
Each Node runs the services needed to support containers.
**Component of Node**
- Kubelet
- Container Engine
- kube-Proxy

1. Kubelet
   - kubelet agent in running on node
   - listen to kubernates master(eg, pod creation request)
   - use port 10255
   - send success/fail reports to master
   - kubelet communicate to controller-manager via api-server
  
2. Container Engine
   - kunernates support multiple engines buy majorly work with docker only.
   - works with kubelet
   - pull image
   - start/stop containers
   - exposing containers on ports specified in manifest
     
3. Kube-proxy
   - Runs on each node.
   - Assigns unique IP addresses to each pod.
   - kube-proxy runs on each node & this make sure that each pod will get its own unique ip address

### Pode
  - Smallest deployable unit in Kubernetes.
  - A Pod contains one or more containers deployed together on the same host.
  - Kubernetes manages pods, not individual containers.
  - One pod per container is standard practice (even though multiple containers are allowed).
  - in kubernates, the controll unit is the pod, not containers
  - pod runs on node which runs by master
  - Cannot start a container without a pod.
    
**multicontainer Pode**
- Containers inside a pod:
  - Share memory/network space.
  - Communicate via localhost.
  - Share the same volume (storage).
  - Deployed in an all-or-nothing manner.


### Higher level kubernates Objects
Replication set: Provides auto-healing and auto-scaling (not provided by default in K8s).\
deployment: Adds versioning and rollback capabilities to ReplicationSets.\
service:	Exposes pods via static IP and DNS, enabling stable networking.\
volume: Provides persistent (non-ephemeral) storage across pod lifecycles.
     
## Custom Installation
Refrance: https://medium.com/@kvihanga/how-to-set-up-a-kubernetes-cluster-on-ubuntu-22-04-lts-433548d9a7d0

### common steps to perform on both master and worker
1. need to create instance with min 2 cpu and 4 gb ram
2. need to allow ports by following: https://kubernetes.io/docs/reference/networking/ports-and-protocols
3. install docker on any container enginer
  ```
  sudo apt install docker.io -y
  ```
4. Install Kubernetes Components
  Add the Kubernetes signing key and repository.
  ```
  sudo apt-get update
  ```
  ```
 sudo apt-get install -y apt-transport-https ca-certificates curl gpg
  ```
  ```
  curl -fsSL https://pkgs.k8s.io/core:/stable:/v1.29/deb/Release.key | sudo gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg
  ```
  ```
  echo 'deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v1.29/deb/ /' | sudo tee /etc/apt/sources.list.d/kubernetes.list
  ```
5. Install kubelet, kubeadm, and kubectl
   ```
   sudo apt-get update
   ```
   ```
   sudo apt-get install -y kubelet kubeadm kubectl
   ```
   ```
   sudo apt-mark hold kubelet kubeadm kubectl
   ```

### steps to perform on master
  1. Install a pod network so that your nodes can communicate with each other
     ```
     kubectl apply -f https://raw.githubusercontent.com/projectcalico/calico/v3.25.0/manifests/calico.yaml
     ```
  2. Initialize Kubernetes Master Node\
     ```
     sudo kubeadm init
     ```

> [!NOTE]
> it will return command at end like kubeadm join.... this command is required to run on pod to connect with it

  3. start using your cluster, set up the kubeconfig
     ```
     mkdir -p $HOME/.kube
     sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
     sudo chown $(id -u):$(id -g) $HOME/.kube/config
     ```
### steps to perform on worker
   1. run command which got end of kubeadm init like this one
      ```
      kubeadm join XX.XXX.XX.XXX:6443 --token z5zbph.rhwgfghx5ymzgyrhqtv \
        --discovery-token-ca-cert-hash sha256:1014dadfgfgfd21eef23caa7dd4caffdbd1185sfdsdffdsfd9a1d99c71e57f25a93dbd1c60c5f12e
      ```


## Commands
- Common commands and their usage.

## Best Practices
- Recommended practices for effective use.

## Security Checks
- Guidelines for ensuring security while using the tool.

## Secret Management
- Strategies for managing sensitive information securely.
