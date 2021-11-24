# func-talk
A brief example on Knative Func using Golang Func Template

First you need to have a Cluster with Knative Serving and an Ingress controller running:
You can find the instructions here:
 [https://github.com/knative-sandbox/kn-plugin-func/blob/main/docs/README.md](https://github.com/knative-sandbox/kn-plugin-func/blob/main/docs/README.md)

Or use the Knative Quickstart Guide here:
[https://knative.dev/docs/getting-started/](https://knative.dev/docs/getting-started/)

 First lets define 2 simples Knative Services (`ksvc`) that are going to be the simulating 2 versions of a long running job, for this lets run:

kn service create long-running --image gcr.io/knative-samples/helloworld-go --port 8080 --env TARGET="First, I've been running for a while"

Then we'll update the `ksvc` so Knative creates for us a new Revision:
```
kn service update long-running --env TARGET="Second, Processing..."
```
Lets confirm that there are currently two Revisions of the `long-running ksvc`
```
kn revision list
```
Now lets tag both Revisions. This because right now the oldest Revisions does not have an url we can use to reach it, with tags, Knative creates an unique Route for each tagged Revision
```
kn service update hello --tag hello-00001=first --tag @latest=second
```
Now is a good time for inspect the fn-golang-example, after all that is the code we are going to be using for this example. The Handle function is the one that have the magic, in summary: we are using a CloudEvent based handler, that call any of our two long running revisions depending on the "service_url" field in it's data body. Then the function returns the response (nothing too fancy but you get the point on things you can do with Knative Func in minutes).
```
func deploy
```

You need to modify the following Routes with your Ksvcs Routes by inspecting the output of:
```
kn ksvc list
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
