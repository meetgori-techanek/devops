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
