go get github.com/facebookincubator/ent/cmd/entc
echo "ent: generate"
entc generate ./database/schema --target ./database/ent/
echo "ent: done"
