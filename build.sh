
mkdir -p 'builds'
mkdir -p 'builds/resources'

oss=(windows darwin linux)
archs=(amd64 arm64)

for os in $oss
do
  mkdir -p 'builds/'"$os"
  for arch in $archs
  do
    env GOOS="$os" GOARCH="$arch" go build -C cmd/main/ -o ../../builds/"$os"/camundaIncidentAggregator_"$arch"
  done
done

cp  -R "./resources/" "./builds/resources"

