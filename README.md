# Func Talk
A brief example on Knative Func using Golang Func Template

## Configuring the Env

 You need to have a Cluster with Knative Serving, an Ingress Controller and The func CLI installed.

 You can find the instructions here:
 [https://github.com/knative-sandbox/kn-plugin-func/blob/main/docs/README.md](https://github.com/knative-sandbox/kn-plugin-func/blob/main/docs/README.md)

 Or use the Knative Quickstart Guide here and then install the `fn CLI`:
[https://knative.dev/docs/getting-started/](https://knative.dev/docs/getting-started/)

## The Example

 First lets define 2 simples Knative Services (`ksvc`) that are going to be the simulating 2 versions of a long running job, for this lets run:
```
kn service create long-running --image gcr.io/knative-samples/helloworld-go --port 8080 --env TARGET="First here, I've been running for a while..."
```

 Then we'll update the `ksvc` so Knative creates for us a new Revision:
```
kn service update long-running --env TARGET="Second here, Processing..."
```

  Lets confirm that there are currently two Revisions of the `long-running ksvc`
```
kn revision list

NAME                 SERVICE        TRAFFIC   TAGS   GENERATION   AGE   CONDITIONS   READY   REASON
long-running-00002   long-running   100%             2            4s    4 OK / 4     True
long-running-00001   long-running                    1            10s   4 OK / 4     True
```

 Now lets add Tags both Revisions. This because right now the oldest Revisions does not have an url we can use to reach it, with Tags, Knative creates an unique Route for each tagged Revision
```
kn service update long-running --tag long-running-00001=first --tag @latest=second
```

 Now is a good time to inspect the fn-golang-example, after all that is the code we are going to be using for this example.
 
 This directory structure was created by running:
```
func create -l go -t cloudevents fn-golang-example
```

 Then we just modify the Handle Functions to do whatever we need.
 The Handle function is the one that have the magic. In summary: we are using a [CloudEvent](https://cloudevents.io/) based handler, that call any of our two `long-running ksvc` revisions depending on the "service_url" field in the CloudEvent data body. This "service_url" contains the url for the specific Tag we want to call.
  Then the function returns the response (nothing too fancy but you get the point on things you can do with Knative Func in minutes).
```
func deploy
```

You need to modify the following Routes with your `ksvc`'s Routes by inspecting the output of:
```
kn ksvc list
NAME           URL                                                    LATEST               AGE     CONDITIONS   READY   REASON
long-running   http://long-running.default.YOUR.IP.IS.HERE.sslip.io   long-running-00002   2m17s   3 OK / 3     True
```

Now lets test this by sending some CloudEvents!

```
curl -X POST \
  -H "content-type: application/json"  \
  -H "ce-specversion: 1.0"  \
  -H "ce-source: firstsource"  \
  -H "ce-type: first.type"  \
  -H "ce-id: 123-abc"  \
  -d '{"service_url":"http://first-long-running.default.X.X.X.X.sslip.io"}' \
  http://fn-golang-example.default.X.X.X.X.sslip.io

Output: Hello First, I've been running for a while!
```
```
curl -X POST \
  -H "content-type: application/json"  \
  -H "ce-specversion: 1.0"  \
  -H "ce-source: firstsource"  \
  -H "ce-type: first.type"  \
  -H "ce-id: 123-abc"  \
  -d '{"service_url":"http://second-long-running.default.X.X.X.X.sslip.io"}' \
  http://fn-golang-example.default.X.X.X.X.sslip.io

Output: Hello Second, Processing...!
```
