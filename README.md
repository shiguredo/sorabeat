# Sorabeat

## カスタム beat 開発

ref.
- Creating a Beat based on Metricbeat | Metricbeat Reference [5.6] | Elastic
  https://www.elastic.co/guide/en/beats/metricbeat/current/creating-beat-from-metricbeat.html

### 準備

- Go
- Python 2.7 (ノ￣￣∇￣￣)ノ‾‾‾━━┻━┻━━
- virtualenv

### 生成

```
go get github.com/elastic/beats/metricbeat
cd $GOPATH/src/github.com/elastic/beats/
git checkout v6.0.0-rc2
python ${GOPATH}/src/github.com/elastic/beats/script/generate.py --type=metricbeat
cd ${GOPATH}/src/github.com/shiguredo/sorabeat
make setup
## module => sora
## metricset => connections
```

### ビルド

```
make
```

or

```
go build -i
```


### 設定ファイル

*注意*

作法として、 `module/sora/_meta/` 以下のファイルを編集するらしい。その後、 `_meta/` 以下と
トップレベル以下の `fields.yml`, `sorabeat.yml` , `sorabeat.reference.yml` が生成される。
トップレベル以下のファイルが使用される。
そのため、トップレベルディレクトリの中の以下のファイルは make ターゲットにより上書きされるので注意。

- fields.yml
- sorabeat.yml
- sorabeat.reference.yml
- modules.d/

*Tips*

- 認証情報、接続情報などを別の YAML ファイルに入れておき、コマンドライン起動時に読み込める
  例: `./sorabeat -c sorabeat.yml -c sorabeat.cred.yml -e -d '*'`

生成

```
make update2
```


### 実行 (debug 用)

```
./sorabeat -c sorabeat.edited.yml -e -d "*"
```

### バージョン設定

```
VERSION=0.1.0 make set_version
```

確認は `make get_version`

*TODO* そのうち git tag と連動したい

### パッケージング

デフォルトでは SNAPSHOT が生成される

```
make package2
```

`build/upload/` 以下にパッケージが生成される。

リリース用

```
SNAPSHOT=false make package2
```

### Linux/ARM64 用バイナリ生成

```
GOOS=linux GOARCH=arm64 make
```

パッケージングは未調査

---------------

以下、生成された README そのまま

sorabeat is a beat based on metricbeat which was generated with metricbeat/metricset generator.


## Getting started

To get started run the following command. This command should only be run once.

```
make setup
```

It will ask you for the module and metricset name. Insert the name accordingly.

To compile your beat run `make`. Then you can run the following command to see the first output:

```
sorabeat -e -d "*"
```

In case further modules are metricsets should be added, run:

```
make create-metricset
```

After updates to the fields or config files, always run

```
make collect
```

This updates all fields and docs with the most recent changes.

## Use vendoring

We recommend to use vendoring for your beat. This means the dependencies are put into your beat folder. The beats team currently uses [govendor](https://github.com/kardianos/govendor) for vendoring.

```
govendor init
govendor update +e
```

This will create a directory `vendor` inside your repository. To make sure all dependencies for the Makefile commands are loaded from the vendor directory, find the following line in your Makefile:

```
ES_BEATS=${GOPATH}/src/github.com/elastic/beats
```

Replace it with:
```
ES_BEATS=./vendor/github.com/elastic/beats
```


## Versioning

We recommend to version your repository with git and make it available on Github so others can also use your project. The initialise the git repository and add the first commits, you can use the following commands:

```
git init
git add README.md CONTRIBUTING.md
git commit -m "Initial commit"
git add LICENSE
git commit -m "Add the LICENSE"
git add .gitignore
git commit -m "Add git settings"
git add .
git reset -- .travis.yml
git commit -m "Add sorabeat"
```

## Packaging

The beat frameworks provides tools to crosscompile and package your beat for different platforms. This requires [docker](https://www.docker.com/) and vendoring as described above. To build packages of your beat, run the following command:

```
make package
```

This will fetch and create all images required for the build process. The hole process to finish can take several minutes.
