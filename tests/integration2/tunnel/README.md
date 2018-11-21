## Running tunnel tests

```
-istio.test.env kubernetes -istio.test.noCleanup -istio.test.kube.config  /Users/istio/Downloads/kubeconfig--istio--istio-dev.yaml  --istio.test.kube.helm.values  global.hub=gcr.io/sap-se-gcp-istio-dev,global.tag=d001323,global.imagePullPolicy=Always
```

