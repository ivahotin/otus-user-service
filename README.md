# otus-user-service
User API

### Директория с чартом сервиса

`infra/user-service`

### Установка

`kubectl create ns monitoring`
`kubectl config set-context --current --namespace=monitoring`
`helm repo add prometheus-community https://prometheus-community.github.io/helm-charts`
`helm upgrade --install -f infra/prometheus/prometheus.yaml prometheus prometheus-community/kube-prometheus-stack --atomic`
`helm upgrade --install -f infra/prometheus/nginx-ingress.yaml nginx ingress-nginx/ingress-nginx --atomic`

`kubectl create ns user-service`
`kubectl config set-context --current --namespace=user-service`
`helm upgrade --install -f infra/user-service/values.yaml user-service infra/user-service/.`

### Тестирование

`newman run integration_tests/user_api_collection.json`

### Удаление

`helm uninstall user-service`
