{{ define "accesslist.yaml.tpl" }}
allowed:
    - spiffe://{{ .Values.global.trustDomain | default "cluster.local" }}/ns/{{ .Release.Namespace }}/sa/istio-mixer-service-account
    - spiffe://{{ .Values.global.trustDomain | default "cluster.local" }}/ns/{{ .Release.Namespace }}/sa/istio-pilot-service-account
{{- end }}
