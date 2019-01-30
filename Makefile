build:
	docker build --no-cache -t alexandreroman/k8s-session-affinity-backend backend && \
	docker build --no-cache -t alexandreroman/k8s-session-affinity-frontend frontend

push: build
	docker push alexandreroman/k8s-session-affinity-backend && \
	docker push alexandreroman/k8s-session-affinity-frontend
