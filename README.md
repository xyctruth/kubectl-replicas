# kubectl-stash

This plugin stash replicas of deployment, save some resources.

## Usage

### Stash replicas of deployment

```bash
$ kubectl-replicas stash -n test # or (kubectl replicas stash -n test)
"app1" stash replicas succeed
"app2" stash replicas succeed
"app3" stash replicas succeed
```

```
$ kubectl get deployments -n test
NAME       READY   STATUS    RESTARTS       AGE
app1       0/1     1            0           141d
app2       0/2     1            0           141d
app3       0/3     1            0           141d
```

### Recover replicas of deployment

```bash
$ kubectl-replicas recover -n test # or (kubectl replicas stash -n test)
"app1" recover replicas 1 succeed
"app2" recover replicas 2 succeed
"app3" recover replicas 3 succeed
```

```bash
$ kubectl get deployments -n test
NAME       READY   STATUS    RESTARTS       AGE
app1       1/1     1            0           141d
app2       2/2     1            0           141d
app3       3/3     1            0           141d
```
