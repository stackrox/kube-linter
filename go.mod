module golang.stackrox.io/kube-linter

go 1.19

require (
	github.com/Masterminds/sprig/v3 v3.2.3
	github.com/bmatcuk/doublestar/v4 v4.6.0
	github.com/cert-manager/cert-manager v1.11.0
	github.com/fatih/color v1.14.1
	github.com/ghodss/yaml v1.0.0
	github.com/golangci/golangci-lint v1.51.2
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/mapstructure v1.5.0
	github.com/openshift/api v3.9.0+incompatible
	github.com/owenrumney/go-sarif/v2 v2.1.2
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.6.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.15.0
	github.com/stretchr/testify v1.8.2
	helm.sh/helm/v3 v3.11.1
	k8s.io/api v0.26.2
	k8s.io/apimachinery v0.26.2
	k8s.io/cli-runtime v0.26.2
	k8s.io/client-go v0.26.2
	k8s.io/gengo v0.0.0-20220902162205-c0856e24416d
)

require (
	4d63.com/gocheckcompilerdirectives v1.2.1 // indirect
	4d63.com/gochecknoglobals v0.2.1 // indirect
	github.com/Abirdcfly/dupword v0.0.9 // indirect
	github.com/Antonboom/errname v0.1.7 // indirect
	github.com/Antonboom/nilnil v0.1.1 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20210617225240-d185dfc1b5a1 // indirect
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/Djarvur/go-err113 v0.0.0-20210108212216-aea10b59be24 // indirect
	github.com/GaijinEntertainment/go-exhaustruct/v2 v2.3.0 // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/Masterminds/semver/v3 v3.2.0 // indirect
	github.com/OpenPeeDeeP/depguard v1.1.1 // indirect
	github.com/alexkohler/prealloc v1.0.0 // indirect
	github.com/alingse/asasalint v0.0.11 // indirect
	github.com/ashanbrown/forbidigo v1.4.0 // indirect
	github.com/ashanbrown/makezero v1.1.1 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bkielbasa/cyclop v1.2.0 // indirect
	github.com/blizzy78/varnamelen v0.8.0 // indirect
	github.com/bombsimon/wsl/v3 v3.4.0 // indirect
	github.com/breml/bidichk v0.2.3 // indirect
	github.com/breml/errchkjson v0.3.0 // indirect
	github.com/butuzov/ireturn v0.1.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/charithe/durationcheck v0.0.9 // indirect
	github.com/chavacava/garif v0.0.0-20221024190013-b3ef35877348 // indirect
	github.com/containerd/containerd v1.6.18 // indirect
	github.com/curioswitch/go-reassign v0.2.0 // indirect
	github.com/cyphar/filepath-securejoin v0.2.3 // indirect
	github.com/daixiang0/gci v0.9.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/denis-tingaikin/go-header v0.4.3 // indirect
	github.com/docker/cli v20.10.21+incompatible // indirect
	github.com/docker/distribution v2.8.1+incompatible // indirect
	github.com/docker/docker v20.10.21+incompatible // indirect
	github.com/docker/docker-credential-helpers v0.7.0 // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-metrics v0.0.1 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/emicklei/go-restful/v3 v3.9.0 // indirect
	github.com/esimonov/ifshort v1.0.4 // indirect
	github.com/ettle/strcase v0.1.1 // indirect
	github.com/evanphx/json-patch v5.6.0+incompatible // indirect
	github.com/fatih/structtag v1.2.0 // indirect
	github.com/firefart/nonamedreturns v1.0.4 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/fzipp/gocyclo v0.6.0 // indirect
	github.com/go-critic/go-critic v0.6.7 // indirect
	github.com/go-errors/errors v1.0.1 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.20.0 // indirect
	github.com/go-openapi/swag v0.19.14 // indirect
	github.com/go-toolsmith/astcast v1.1.0 // indirect
	github.com/go-toolsmith/astcopy v1.0.3 // indirect
	github.com/go-toolsmith/astequal v1.1.0 // indirect
	github.com/go-toolsmith/astfmt v1.1.0 // indirect
	github.com/go-toolsmith/astp v1.1.0 // indirect
	github.com/go-toolsmith/strparse v1.1.0 // indirect
	github.com/go-toolsmith/typep v1.1.0 // indirect
	github.com/go-xmlfmt/xmlfmt v1.1.2 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/gofrs/flock v0.8.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golangci/check v0.0.0-20180506172741-cfe4005ccda2 // indirect
	github.com/golangci/dupl v0.0.0-20180902072040-3e9179ac440a // indirect
	github.com/golangci/go-misc v0.0.0-20220329215616-d24fe342adfe // indirect
	github.com/golangci/gofmt v0.0.0-20220901101216-f2edd75033f2 // indirect
	github.com/golangci/lint-1 v0.0.0-20191013205115-297bf364a8e0 // indirect
	github.com/golangci/maligned v0.0.0-20180506175553-b1d89398deca // indirect
	github.com/golangci/misspell v0.4.0 // indirect
	github.com/golangci/revgrep v0.0.0-20220804021717-745bb2f7c2e6 // indirect
	github.com/golangci/unconvert v0.0.0-20180507085042-28b1c447d1f4 // indirect
	github.com/google/btree v1.0.1 // indirect
	github.com/google/gnostic v0.6.9 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/gordonklaus/ineffassign v0.0.0-20230107090616-13ace0543b28 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/gostaticanalysis/analysisutil v0.7.1 // indirect
	github.com/gostaticanalysis/comment v1.4.2 // indirect
	github.com/gostaticanalysis/forcetypeassert v0.1.0 // indirect
	github.com/gostaticanalysis/nilerr v0.1.1 // indirect
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-version v1.6.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hexops/gotextdiff v1.0.3 // indirect
	github.com/huandu/xstrings v1.3.3 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/inconshreveable/mousetrap v1.0.1 // indirect
	github.com/jgautheron/goconst v1.5.1 // indirect
	github.com/jingyugao/rowserrcheck v1.1.1 // indirect
	github.com/jirfag/go-printf-func-name v0.0.0-20200119135958-7558a9eaa5af // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/julz/importas v0.1.0 // indirect
	github.com/junk1tm/musttag v0.4.5 // indirect
	github.com/kisielk/errcheck v1.6.3 // indirect
	github.com/kisielk/gotool v1.0.0 // indirect
	github.com/kkHAIKE/contextcheck v1.1.3 // indirect
	github.com/klauspost/compress v1.14.1 // indirect
	github.com/kulti/thelper v0.6.3 // indirect
	github.com/kunwardeep/paralleltest v1.0.6 // indirect
	github.com/kyoh86/exportloopref v0.1.11 // indirect
	github.com/ldez/gomoddirectives v0.2.3 // indirect
	github.com/ldez/tagliatelle v0.4.0 // indirect
	github.com/leonklingele/grouper v1.1.1 // indirect
	github.com/liggitt/tabwriter v0.0.0-20181228230101-89fcab3d43de // indirect
	github.com/lufeee/execinquery v1.2.1 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mailru/easyjson v0.7.6 // indirect
	github.com/maratori/testableexamples v1.0.0 // indirect
	github.com/maratori/testpackage v1.1.0 // indirect
	github.com/matoous/godox v0.0.0-20210227103229-6504466cf951 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/mbilski/exhaustivestruct v1.2.0 // indirect
	github.com/mgechev/revive v1.2.5 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/moby/locker v1.0.1 // indirect
	github.com/moby/term v0.0.0-20221205130635-1aeaba878587 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/monochromegane/go-gitignore v0.0.0-20200626010858-205db1a8cc00 // indirect
	github.com/moricho/tparallel v0.2.1 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/nakabonne/nestif v0.3.1 // indirect
	github.com/nbutton23/zxcvbn-go v0.0.0-20210217022336-fa2cb2858354 // indirect
	github.com/nishanths/exhaustive v0.9.5 // indirect
	github.com/nishanths/predeclared v0.2.2 // indirect
	github.com/nunnatsa/ginkgolinter v0.8.1 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0-rc2 // indirect
	github.com/pelletier/go-toml/v2 v2.0.6 // indirect
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/polyfloyd/go-errorlint v1.1.0 // indirect
	github.com/prometheus/client_golang v1.14.0 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.37.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	github.com/quasilyte/go-ruleguard v0.3.19 // indirect
	github.com/quasilyte/gogrep v0.5.0 // indirect
	github.com/quasilyte/regex/syntax v0.0.0-20200407221936-30656e2c4a95 // indirect
	github.com/quasilyte/stdinfo v0.0.0-20220114132959-f7386bf02567 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/ryancurrah/gomodguard v1.3.0 // indirect
	github.com/ryanrolds/sqlclosecheck v0.4.0 // indirect
	github.com/sanposhiho/wastedassign/v2 v2.0.7 // indirect
	github.com/sashamelentyev/interfacebloat v1.1.0 // indirect
	github.com/sashamelentyev/usestdlibvars v1.23.0 // indirect
	github.com/securego/gosec/v2 v2.15.0 // indirect
	github.com/shazow/go-diff v0.0.0-20160112020656-b6b7b6733b8c // indirect
	github.com/shopspring/decimal v1.2.0 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	github.com/sivchari/containedctx v1.0.2 // indirect
	github.com/sivchari/nosnakecase v1.7.0 // indirect
	github.com/sivchari/tenv v1.7.1 // indirect
	github.com/sonatard/noctx v0.0.1 // indirect
	github.com/sourcegraph/go-diff v0.7.0 // indirect
	github.com/spf13/afero v1.9.3 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/ssgreg/nlreturn/v2 v2.2.1 // indirect
	github.com/stbenjam/no-sprintf-host-port v0.1.1 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	github.com/subosito/gotenv v1.4.2 // indirect
	github.com/t-yuki/gocover-cobertura v0.0.0-20180217150009-aaee18c8195c // indirect
	github.com/tdakkota/asciicheck v0.1.1 // indirect
	github.com/tetafro/godot v1.4.11 // indirect
	github.com/timakin/bodyclose v0.0.0-20221125081123-e39cf3fc478e // indirect
	github.com/timonwong/loggercheck v0.9.3 // indirect
	github.com/tomarrell/wrapcheck/v2 v2.8.0 // indirect
	github.com/tommy-muehle/go-mnd/v2 v2.5.1 // indirect
	github.com/ultraware/funlen v0.0.3 // indirect
	github.com/ultraware/whitespace v0.0.5 // indirect
	github.com/uudashr/gocognit v1.0.6 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20180127040702-4e3ac2762d5f // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v1.2.0 // indirect
	github.com/xlab/treeprint v1.1.0 // indirect
	github.com/yagipy/maintidx v1.0.0 // indirect
	github.com/yeya24/promlinter v0.2.0 // indirect
	gitlab.com/bosi/decorder v0.2.3 // indirect
	go.starlark.net v0.0.0-20200306205701-8dd3e2ee1dd5 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/goleak v1.1.12 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	go.uber.org/zap v1.24.0 // indirect
	golang.org/x/crypto v0.5.0 // indirect
	golang.org/x/exp v0.0.0-20220722155223-a9213eeb770e // indirect
	golang.org/x/exp/typeparams v0.0.0-20230203172020-98cc5a0785f9 // indirect
	golang.org/x/mod v0.8.0 // indirect
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/oauth2 v0.4.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/term v0.5.0 // indirect
	golang.org/x/text v0.7.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	golang.org/x/tools v0.6.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20221227171554-f9683d7f8bef // indirect
	google.golang.org/grpc v1.52.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	honnef.co/go/tools v0.4.2 // indirect
	k8s.io/apiextensions-apiserver v0.26.0 // indirect
	k8s.io/klog/v2 v2.80.1 // indirect
	k8s.io/kube-openapi v0.0.0-20221207184640-f3cff1453715 // indirect
	k8s.io/utils v0.0.0-20221128185143-99ec85e7a448 // indirect
	mvdan.cc/gofumpt v0.4.0 // indirect
	mvdan.cc/interfacer v0.0.0-20180901003855-c20040233aed // indirect
	mvdan.cc/lint v0.0.0-20170908181259-adc824a0674b // indirect
	mvdan.cc/unparam v0.0.0-20221223090309-7455f1af531d // indirect
	oras.land/oras-go v1.2.2 // indirect
	sigs.k8s.io/gateway-api v0.6.0 // indirect
	sigs.k8s.io/json v0.0.0-20220713155537-f223a00ba0e2 // indirect
	sigs.k8s.io/kustomize/api v0.12.1 // indirect
	sigs.k8s.io/kustomize/kyaml v0.13.9 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)
