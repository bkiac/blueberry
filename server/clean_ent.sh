echo "ent: clean"
find ./database/ent -mindepth 1 ! -regex '^./database/ent/README.md' -delete
