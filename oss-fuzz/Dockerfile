FROM gcr.io/oss-fuzz-base/base-builder-go

COPY . $SRC/kubearmor

COPY ./oss-fuzz/build.sh $SRC/

ENV FUZZING_LANGUAGE=go

WORKDIR $SRC/kubearmor/KubeArmor
