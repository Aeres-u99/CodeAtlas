package internal

type Call struct {
	Caller string
	Callee string
}

type Dependency struct {
	Source string
	Target string
}

const DefaultHermesIgnore = `.git/
.gitmodules
.gitignore
.gitattributes
.gitkeep
bin/
build/
dist/
out/
target/
release/
debug/
obj/
objs/
tmp/
temp/
*.o
*.a
*.so
*.dll
*.dylib
*.exe
*.out
*.class
*.jar
*.war
*.ear
vendor/
pkg/
node_modules/
bower_components/
__pycache__/
*.pyc
*.pyo
*.pyd
.venv/
venv/
env/
.python-version
coverage.out
cover.out
*.test
target/
.next/
.nuxt/
.svelte-kit/
.parcel-cache/
.turbo/
.eslintcache
*.tsbuildinfo
.idea/
.vscode/
.vs/
*.swp
*.swo
*~
*.log
logs/
.DS_Store
Thumbs.db
desktop.ini
.cache/
.cache-loader/
.ccache/
coverage/
htmlcov/
*.prof
*.pprof
*.gcda
*.gcno
*.gcov
site/
docs/_build/
*.zip
*.tar
*.tar.gz
*.tgz
*.7z
*.rar
.env
.env.*
*.pem
*.key
*.crt
*.p12
*.db
*.sqlite
*.sqlite3
.terraform/
*.tfstate
*.tfstate.*
*.kubeconfig
hermes.json
`
