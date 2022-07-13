apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: binpak-backend
  name: binpak-backend
  namespace: binpak
spec:
  replicas: 1
  selector:
    matchLabels:
      app: binpak-backend
  template:
    metadata:
      labels:
        app: binpak-backend
    spec:
      serviceAccountName: binpak-backend
      containers:
      - image: ${BACKEND_IMAGE_TAG}
        name: binpak-backend
        ports:
        - containerPort: 8080
          name: http
        resources: {}
        imagePullPolicy: Always
