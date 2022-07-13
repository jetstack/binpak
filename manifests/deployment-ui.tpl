apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: binpak-ui
  name: binpak-ui
  namespace: binpak
spec:
  replicas: 1
  selector:
    matchLabels:
      app: binpak-ui
  template:
    metadata:
      labels:
        app: binpak-ui
    spec:
      serviceAccountName: binpak-ui
      containers:
      - image: ${UI_IMAGE_TAG}
        name: binpak-ui
        ports:
        - containerPort: 80
          name: http
        resources: {}
        imagePullPolicy: Always
