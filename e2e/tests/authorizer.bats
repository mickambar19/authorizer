#!/usr/bin/env bats
setup() {
  load '../node_modules/bats-support/load'
  load '../node_modules/bats-assert/load'
  go build -o ./authorizer ../src/main.go
}

setup

helper() {
  folderName=$1
  echo $folderName
  inputFileName="./cases/${folderName}/input"
  ouputFileName="./cases/${folderName}/output"
  echo $inputFileName
  echo $ouputFileName
  actual="$(./authorizer <$inputFileName)"
  expectedOutput="$(cat <$ouputFileName)"
  assert_equal "$actual" "$expectedOutput"
}

@test "Account-Initialize" {
  helper "account-initialize"
}

@test "Account Initialize Violaton" {
  helper "account-initialize-violaton"
}

@test "Transaction Account Not Initialized" {
  helper "transaction-account-not-initialized"
}

@test "Transaction Card Not Active" {
  helper "transaction-card-not-active"
}

@test "Transaction Double Transaction" {
  helper "transaction-double-transaction"
}

@test "Transaction High Frequency" {
  helper "transaction-high-frequency"
}

@test "Transaction insufficient-limit" {
  helper "transaction-insufficient-limit"
}

@test "Transaction Multiple Violations" {
  helper "transaction-multiple-violations"
}

@test "Transaction Sucessfully" {
  helper "transaction-sucessfully"
}
