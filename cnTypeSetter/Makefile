.RECIPEPREFIX = >

all: linux_x64 linux_arm

linux_x64:
>  go build -o ${GOBIN}/cnTypeSetterAMD -v .

# Run on a Raspberry PI (raspbian)
linux_arm:
>  env GOARCH=arm go build -o ${GOBIN}/cnTypeSetterARM -v .
