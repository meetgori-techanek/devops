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
###Selectors
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
###Benefits:
1. Reliablity: by having multiple versions of application will prevent outage of aaplication if one or more pods failed
2. load balancing: having multiple version of containers enables you to easily send traffic instances to prevent overloading of a single instances or node
3. scaling: when load does become too much for the number of existing instances, kubernates enables you to easily scale up application
4. rolling updates: updates to a service by replacing pods one by one

### Replication Controller


