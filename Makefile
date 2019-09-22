define install_istio =
	ISTIO_PATH=/home/augusto.pimenta/Downloads/istio-1.3.0/
	for i in ${ISTIO_PATH}/install/kubernetes/helm/istio-init/files/crd*yaml; do kubectl apply -f $i; done
	kubectl apply -f ${ISTIO_PATH}/install/kubernetes/istio-demo.yaml
	kubectl get svc -n istio-system
	kubectl get pods -n istio-system
endef

define install_kiali_ui =
	ISTIO_PATH=/home/augusto.pimenta/Downloads/istio-1.3.0/
	KIALI_USERNAME=$(echo -n 'admin' | base64)
	KIALI_PASSPHRASE=$(echo -n 'admin' | base64)
	kubectl apply -f kiali-secret.yml
	helm template \
	--set kiali.enabled=true \
	--set "kiali.dashboard.jaegerURL=http://jaeger-query:16686" \
	--set "kiali.dashboard.grafanaURL=http://grafana:3000" \
	${ISTIO_PATH}/install/kubernetes/helm/istio \
	--name istio --namespace istio-system > $HOME/istio.yaml
	kubectl apply -f $HOME/istio.yaml
endef

define start_grafana =
	kubectl -n istio-system port-forward $(kubectl -n istio-system get pod -l app=grafana -o jsonpath='{.items[0].metadata.name}') 3000:3000
endef

define start_kiali_ui =
	kubectl -n istio-system port-forward $(kubectl -n istio-system get pod -l app=kiali -o jsonpath='{.items[0].metadata.name}') 20001:20001
endef

start-minikube:
	minikube start --memory=8192 --cpus=4 --kubernetes-version=v1.14.2
delete-minikube:
	minikube delete
tunnel-minikube:
	minikube tunnel --cleanup	
	minikube tunnel
grafana: ; $(value start_grafana)

kiali: ; $(value start_kiali_ui)

install-istio: ; $(value install_istio)

install-kiali-ui: ; $(value install_kiali_ui)

install-raptorslog:
	kubectl label namespace default istio-injection=enabled
	kubectl apply -f k8s/

install-all: 
	start-minikube
	install-istio
	install-kiali-ui

.ONESHELL: