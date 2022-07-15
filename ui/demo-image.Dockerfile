# Create a combined image that serves both the UI and a mock JSON API response
FROM europe-west1-docker.pkg.dev/jetstack-wil/binpak/binpak-ui:latest

COPY demo-image-nginx.conf /etc/nginx/conf.d/default.conf
COPY demo-image-api-response.json /backend-mock.json