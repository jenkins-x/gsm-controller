module github.com/REPLACE_ME_ORG/REPLACE_ME_APP_NAME

go 1.12

require (
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.6.2
	github.com/stretchr/testify v1.4.0
)

replace github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.0.0

replace github.com/russross/blackfriday => github.com/russross/blackfriday v1.5.2
