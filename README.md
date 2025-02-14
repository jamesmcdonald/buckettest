# buckettest

A dumb little web tool to test access to a GCS bucket. You can easily destroy
data with this, so don't use it on important buckets.

It doesn't try to do any authentication, it just relies on automatic
credentials. Locally that'll be gcloud's application default credentials, but
the point of this is to test Workload Identity.

With Workload Identity Federation you shouldn't need to configure anything
special, just use the service account that has been granted bucket access. With
regular Workload Identity you'll need an annotation on the service account to
map to the correct GCP service account. That looks like:

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: buckettest
  annotations:
    iam.gke.io/gcp-service-account: sa-name@gcp-project.iam.gserviceaccount.com
```
