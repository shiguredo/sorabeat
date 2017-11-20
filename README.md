[![CircleCI](https://circleci.com/gh/shiguredo/sorabeat/tree/develop.svg?style=svg)](https://circleci.com/gh/shiguredo/sorabeat/tree/develop)

# Sorabeat

Sorabeat は [WebRTC SFU Sora](https://sora.shiguredo.jp) の統計情報を Elasticsearch  や Logstash
に送信するソフトウェアです。

Elastic Beats のひとつとして、libbeat, metricbeat を基に作られています。

## インストール

releases から RPM/DEB/tar.gz を取得してインストールします。

## 設定

`/etc/sorabeat/sorabeat.yml`, `/etc/sorabeat/modules.d/sora.yml` を編集します。

`sorabeat.yml` では、Elasticsearch などの接続先や、認証情報を設定します。
一般的な設定は metricbeat と同様です。
[Configuring Metricbeat | Metricbeat Reference \[6.0\] | Elastic](https://www.elastic.co/guide/en/beats/metricbeat/6.0/configuring-howto-metricbeat.html)
を参照してください。

`sora.yml` では、Sora の接続先、データ取得間隔の設定を行います。
以下にサンプルを示します。

```
- module: sora
  metricsets: ["stats", "connections"]
  # 更新間隔を指定します
  period: 5s
  # Sora サーバの API アクセスのホスト、ポートを指定します。
  hosts: ["127.0.0.1:3000"]
```

## 起動

RPM でインストールした場合、service コマンドで起動、終了を制御できます。

```
service sorabeat start
```

ログは `/var/log/sorabeat/` 以下に出力されます。

*TODO* : DEB, tar.gz インストールのときの使い方を追加する

## Elasticsearch インデックス

Elasticsearch のインデックスパターンは、 `sorabeat-*` です。

## stats メトリックセット

ソースは Sora の `GetStatsReport` です。
フィールド名は `sora.stats.` をプレフィックスに持ちます。例えば `average_duration_sec`
は Elasticsearch では `sora.stats.average_duration_sec` フィールドに対応します。

### Sorabeat が追加するフィールド

以下の項目はもともと数値のリストです。

- erlang_vm.statistics.active_tasks
- erlang_vm.statistics.active_tasks_all
- erlang_vm.statistics.run_queue_lengths
- erlang_vm.statistics.run_queue_lengths_all

それらに対して最大値(max)、最小値(min)、平均値(mean)、標準偏差(stdddv)と
不均衡さ(imbalance)を
フィールドとして追加します。`erlang_vm.statistics.active_tasks` を例に取ると
次のフィールドが追加されます。

- sora.stats.erlang_vm.statistics.active_tasks_max
- sora.stats.erlang_vm.statistics.active_tasks_min
- sora.stats.erlang_vm.statistics.active_tasks_mean
- sora.stats.erlang_vm.statistics.active_tasks_stddev
- sora.stats.erlang_vm.statistics.active_tasks_imbalance

各リスト values の不均衡さ(imbalance)は次で計算しています。

```
                    最大値(values)
不均衡さ = ------------------------------
             最大値(最小値(values) , 1)
```


## connections メトリックセット

ソースは Sora の `GetStatsAllConnections` です。
フィールド名は `sora.connections.` をプレフィックスに持ちます。例えば `rtp` の下にある
`total_received_bytes` は Elasticsearch では `sora.connections.rtp.total_received_bytes`
フィールドに対応します。

### Sorabeat が追加するフィールド

- `sora.connections.channel_client_id`: `channel_id` と `client_id` を
  スラッシュ (`/`) で結合した文字列

## dashboard, visualization のセットアップ

`sorabeat setup` を実行すると各数値型フィールドの visualization とサンプルの簡単なダッシュボードが
ロードされます。
適切な権限をもったユーザと、kibana の endpoint 設定が必要です。


------------

# 以下、開発者向け

## カスタム beat 開発の参考

- Beats 開発全般

  - Beats Developer Guide [master] | Elastic
    https://www.elastic.co/guide/en/beats/devguide/current/index.html
- Metricbeat をベースにしたカスタム beat 開発 (Sorabeat はコレ)

  - Creating a Beat based on Metricbeat | Beats Developer Guide [master] | Elastic
    https://www.elastic.co/guide/en/beats/devguide/current/creating-beat-from-metricbeat.html

- Beat や Beat module のための Kibana ダッシュボードを作る方法

  - Creating New Kibana Dashboards for a Beat or a Beat module | Beats Developer Guide [master] | Elastic
    https://www.elastic.co/guide/en/beats/devguide/current/new-dashboards.html

- 以下、Sorabeat には関係ないが、近隣なので参考まで

  - Metricbeat のモジュールだけを新規で開発

    - Creating a Metricbeat Module | Beats Developer Guide [master] | Elastic
      https://www.elastic.co/guide/en/beats/devguide/current/creating-metricbeat-module.html

  - イチからカスタム Beat を開発 (Sorabeat には関係ない、参考まで)

    - Creating a New Beat | Beats Developer Guide [master] | Elastic
      https://www.elastic.co/guide/en/beats/devguide/current/new-beat.html

### 準備

- Go 1.9.2
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
## 対話形式で進むので入力
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

これらは最終的にパッケージに入るので、認証情報を入れないこと。

*Tips*

- 認証情報、接続情報などを別の YAML ファイルに入れておき、コマンドライン起動時に読み込める。
  複数も可能。例: `./sorabeat -c sorabeat.yml -c sorabeat.cred.yml -e -d '*'`

生成

```
make update2
```

*NOTE* 生成されるファイルが metricbeat となる(libbeat/metricbeat での抽象化不足?)バグのため、
update target を update2 で少々上書きしている。以下に出てくる set_version2, package2 も同様。

*TODO* ↑の issue を beats repogitory に切る

### 実行 (debug 用)

```
./sorabeat -c sorabeat.edited.yml -e -d "*"
```

起動オプション

| オプション      | 説明                                 |
|-----------------|--------------------------------------|
| -e              | ログをファイルではなく stderr に出す |
| -d \<selector\> | デバッグセレクタを有効にする         |
|                 | セレクタはコード読むしかなさそう     |

セレクタ

- `cfgfile` : 設定ファイルまわり
- `publish` : es / logstash への送信
- `modules` : 読み込まれたモジュールを羅列 (metricbeat.go)

### バージョン設定

```
VERSION=0.1.0 make set_version2
```

確認は `make get_version`

*TODO* そのうち git tag と連動したい

### パッケージ生成

デフォルトでは SNAPSHOT が生成される

```
make clean
make python-env
make package2
```

`build/upload/` 以下にパッケージが生成される。

リリース用

```
make clean
make python-env
SNAPSHOT=false make package2
```

### Linux/ARM64 用バイナリ生成

```
GOOS=linux GOARCH=arm64 make
```

ARM 向けパッケージングは未調査

### beats のバージョン更新

```
cd /path/to/sorabeat
rm -rf vendor
cd $GOPATH/src/github.com/elastic/beats
git fetch
git checkout v6.0.0 # バージョン指定すること
make copy-vendor
```

### visualization / dashboard の生成, fields.yml の生成

単純な visualization をスクリプト `scripts/visualization_single.sh` で生成している。
入力が `scripts/sora_fields.yml` で、出力が `_meta/kibana/default/dashboard/sorabeat_vis1.json` である。

### 手で作った dashboard の保存

Kibana で dashboard の ID (`28516270-bec0-11e7-b277-79c0643bd2c8` のような文字列)を確認して、
curl で Kibana API を叩くと取れる。

例

```
KIBANA_BASE=https://foo.example.com
USER_CRED='kibana_user:its_password'

curl -s \
     ${KIBANA_BASE}/api/kibana/dashboards/export'?dashboard='${DASHBOARD} \
     -u ${USER_CRED}

```

dashboard としては練習で作ったものを export したものを
`_meta/kibana/default/dashboard/exported-dashboard1.json` に入れている。

## TODO

- fields.yml も sora_fields.yml から生成できるようにしたい
- dashboard を充実させる
- fields.yml に無駄なフィールドが入っている
- visualization で hostname フィルタ(クエリ)が Kibana UI として書けないか調べる
- ARM64 パッケージング
- パッケージを絞ってビルドを早くする

