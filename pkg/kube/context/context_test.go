package context

var (
	clusterFooContext = "clustuer-foo"
	clusterBarContext = "clustuer-bar"
)

const validKubeConfig = `
apiVersion: v1
kind: Config
clusters:
- cluster:
    server: https://foo.com
  name: cluster-foo
- cluster:
    server: https://bar.com
  name: cluster-bar
contexts:
- context:
    cluster: cluster-foo
    user: user1
  name: cluster-foo
- context:
    cluster: cluster-bar
    user: user1
  name: cluster-bar
current-context: cluster-foo
users:
- name: user1
  user:
    password: secret
    username: user
`

const changeKubeConfig = `
apiVersion: v1
kind: Config
clusters:
- cluster:
    server: https://changed-url.com
  name: cluster-bar
contexts:
- context:
    cluster: cluster-bar
    user: user-bar
  name: context-bar
- context:
    cluster: cluster-bar
    user: user-baz
  name: context-baz
current-context: context-baz
users:
- name: user1
  user:
    password: secret
    username: user
`
