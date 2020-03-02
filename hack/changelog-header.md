### Linux

```shell
curl -L https://github.com/jenkins-x-labs/gsm-controller/releases/download/v{{.Version}}/gsm-controller-linux-amd64.tar.gz | tar xzv 
sudo mv gsm-controller /usr/local/bin
```

### macOS

```shell
curl -L  https://github.com/jenkins-x-labs/gsm-controller/releases/download/v{{.Version}}/gsm-controller-darwin-amd64.tar.gz | tar xzv
sudo mv gsm-controller /usr/local/bin
```

