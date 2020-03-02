### Linux

```shell
curl -L https://github.com/REPLACE_ME_ORG/REPLACE_ME_APP_NAME/releases/download/v{{.Version}}/REPLACE_ME_APP_NAME-linux-amd64.tar.gz | tar xzv 
sudo mv REPLACE_ME_APP_NAME /usr/local/bin
```

### macOS

```shell
curl -L  https://github.com/REPLACE_ME_ORG/REPLACE_ME_APP_NAME/releases/download/v{{.Version}}/REPLACE_ME_APP_NAME-darwin-amd64.tar.gz | tar xzv
sudo mv REPLACE_ME_APP_NAME /usr/local/bin
```

