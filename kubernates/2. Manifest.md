# manifest components

## Labels and Selectors
### Labels 
```
..........
  labels:                                                   
    env: development
    class: pods
..........
```
```
kuberctl label pods <pod-name> key=value
```

- labels are mechanism used to organise kubernates objects
- it is in key-value format
- it can be refer to development environment like dev, stage, prod or product group like deptA,deptB etc

**To see lable of pod/node**
```
kubectl get pods/nodes --show-labels
```
### Selectors
**List/delete pod with specific label**

- Equality based label(=,!=)
```
kubectl get/delete pod -l env=development
```

- Set based label(in, not in, exist)\
     < key > in(val1,val2)
```
kubectl get/delete pod -l 'env in development'
```

**node Selectors**
```
..........
    nodeSelector:                                         
       hardware: test-pod
..........
```
then we have to manually tag node with label 

```
kuberctl label pods <pod-name> hardware=test-pod
```

- one case for selecting labels it to constrain the set of nodes into which a pod can schedule
- generally such constraints are not happening, as the scheduler will automatically do a reasonable placement of pods, but in certain circumstances we might need it.
- first we give label to node, then use node selector tot the pod configuration

## Scalling and replication
- kubernates was designed to orchestate multiple containers and replications
### Benefits:
1. Reliablity: by having multiple versions of application will prevent outage of aaplication if one or more pods failed
2. load balancing: having multiple version of containers enables you to easily send traffic instances to prevent overloading of a single instances or node
3. scaling: when load does become too much for the number of existing instances, kubernates enables you to easily scale up application
4. rolling updates: updates to a service by replacing pods one by one

### Replication Controller
- A replication controller is a object that enables you to easily create multiple pods then make sure that numbers of pods always exist
- if a pod created using replication controller will be automatically replaced if they crash, failed, or terminated
- replication controller is recomended if you just want to make sure atleast 1 pod is always running even after systen reboot
   
**Sample**
```
kind: ReplicationController
apiVersion: v1
metadata:
  name: myreplica
spec:
  replicas: 2
  selector:
    myname: meet
  template:
    metadata:
      name: testpod6
      labels:
        myname: meet
    spec:
     containers:
       - name: c00
         image: ubuntu
         command: ["/bin/bash", "-c", "while true; do echo Hello-Meet; sleep 5 ; done"]
```

- all replication controller running
```
 kubectl get rc
```

- all details of replica controller
```
 kubectl describe rc <replica-controller-name>
```
- scale replicas with command
```
 kubectl scale --replicas=8 rc -l key=value
```
- delete pods
```
 kubectl delete -f <file.yml>
```

### Replicaset
- replica set is next version of replication controller
- replication controller is only equality based(=,!=) selector, where as the replica set supports set-based collector
- replicaset rather than the replication controller is used by other objects like deployment

  **Sample**
```
kind: ReplicaSet
apiVersion: apps/v1
metadata:
  name: myrs
spec:
  replicas: 2
  selector:
    matchExpressions:                             # these must match the labels
       - {key: myname, operator: In, values: [meet, mit]}
       - {key: env, operator: NotIn, values: [production]}
  template:
    metadata:
      name: testpod7
      labels:
        myname: meet
    spec:
     containers:
       - name: c00
         image: ubuntu
         command: ["/bin/bash", "-c", "while true; do echo Technical-Guftgu; sleep 5 ; done"]
```

- all replicationset running
```
 kubectl get rs
```

- all details of replicaset
```
 kubectl describe rs <replica-controller-name>
```
- scale replicas with command
```
 kubectl scale --replicas=8rsrc -l key=value
```
- delete pods
```
 kubectl delete -f <file.yml>
```

### Deployment & Rollback
Replication controller and replicaset is not able to do update & rollback apps in the cluster 

- deployment object act as a suprevisor for pods, giving you fine-grained control over how and when a new pod is rolled out, updated or rollback to previous state
- when using deployment object, we first define the state of of the app, then k8s cluster schedules mentioned app instance onto specific individual nodes
- k8s then monitors, if the node hosting on instance goes down or pod is deleted the deplotment controller replaces it.
- this provides a self healing mechanism to address machine failure or maintainance
- **A deplotment provides declarative updates for pods and replicas**

**Typical use cases of deployments**
1. create a deployment to rollout a replicaset -> the replicaset creates pods in the background. check the status of the rollout to see if it succeeds or not.
2. declare the new state of the pods -> by updating the podtemplatespec of the deployment. A new replicaset is created and the deployment manages moving the pods from the old replicaset to new one at a controlled rate. Each new replicaset  updates the revision of the deployment.
3. rollback to earlier deployment revision -> if the current state of the deployment is not stable. each rollback updates revision to earlier deplyment
4. scale up the deployment to facilitates more load
5. pause the deployment to apply multiple fixes to its podtemplatespec and then resume it to start a new rollout.
6. cleanup older replicaset that you don't need anymore 

**Reasons for failed deployments**
1. insufficiant quota
2. readiness probe failures
3. image pull errors
4. insufficiant permissions
5. limit ranges
6. application runtime misconfigurations 

**sample**
```
kind: Deployment
apiVersion: apps/v1
metadata:
   name: mydeployments
spec:
   replicas: 2
   selector:     
    matchLabels:
     name: deployment
   template:
     metadata:
       name: testpod
       labels:
         name: deployment
     spec:
      containers:
        - name: c00
          image: ubuntu
          command: ["/bin/bash", "-c", "while true; do echo Technical-Guftgu; sleep 5; done"]
```

- list all deployments
```
 kubectl get deploy
```

- all details of replicaset
```
 kubectl describe deploy <replica-controller-name>
```
- scale replicas with command
```
 kubectl scale --replicas=8 deploy <deployment-name>
```
- delete pods
```
 kubectl delete -f <file.yml>
```
- get logs from pod
```
 kubectl get logs -f <pod-name>
```
---

 - get status of deployments
```
kubectl rollout status deployment <deployment-name>
```
- get history ot rollouts
```
kubectl rollout history deployment <deployment-name>
```
- rollback to previous version
```
kubectl rollout undo deploy/<deployment-name>
```
