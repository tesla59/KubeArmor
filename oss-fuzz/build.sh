#!/bin/bash -eu
# Dependencies for compile_native_go_fuzzer
go install github.com/AdamKorcz/go-118-fuzz-build@latest
go get github.com/AdamKorcz/go-118-fuzz-build/testing

go get github.com/kubearmor/KubeArmor/KubeArmor/feeder

compile_native_go_fuzzer github.com/kubearmor/KubeArmor/KubeArmor/feeder FuzzU fuzz_u
compile_native_go_fuzzer github.com/kubearmor/KubeArmor/KubeArmor/feeder FuzzFeeder_PushLog fuzz_feeder_push_log
compile_native_go_fuzzer github.com/kubearmor/KubeArmor/KubeArmor/monitor FuzzUpdateLogs fuzz_update_logs
