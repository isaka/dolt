#!/usr/bin/env bats
load $BATS_TEST_DIRNAME/helper/common.bash

setup() {
    setup_common

    dolt sql <<SQL
CREATE TABLE test (
    pk int primary key
);

INSERT INTO test VALUES (0),(1),(2);
SQL
}

teardown() {
    assert_feature_version
    teardown_common
}

@test "sql-branch: DOLT_BRANCH works" {
    run dolt branch
    [[ ! "$output" =~ "new_branch" ]] || false

    run dolt sql -q "SELECT DOLT_BRANCH('new-branch')"
    [ $status -eq 0 ]

    # dolt sql -q "select dolt_branch() should not change the branch
    # It changes the branch for that session which ends after the SQL
    # statements are executed.
    run dolt status
    [ $status -eq 0 ]
    [[ "$output" =~ "main" ]] || false

    run dolt branch
    [ $status -eq 0 ]
    [[ "$output" =~ "new-branch" ]] || false
}

@test "sql-branch: DOLT_BRANCH throws error" {
    # branches that already exist
    dolt branch existing_branch
    run dolt sql -q "SELECT DOLT_BRANCH('existing_branch')"
    [ $status -eq 1 ]
    [ "$output" = "fatal: A branch named 'existing_branch' already exists." ]

    # empty branch
    run dolt sql -q "SELECT DOLT_BRANCH('')"
    [ $status -eq 1 ]
    [ "$output" = "error: cannot branch empty string" ]
}

@test "sql-branch: DOLT_BRANCH -c copies not current branch and stays on current branch" {
    dolt add . && dolt commit -m "0, 1, and 2 in test table"
    run dolt status
    [[ "$output" =~ "main" ]] || false

    dolt checkout -b original
    dolt sql -q "insert into test values (4);"
    dolt add .
    dolt commit -m "add 4 in original"

    dolt checkout main

    # Current branch should be still main with test table without entry 4
    run dolt sql << SQL
SELECT DOLT_BRANCH('-c', 'original', 'copy');
SELECT * FROM test WHERE pk > 3;
SQL
    [ $status -eq 0 ]
    [[ ! "$output" =~ "4" ]] || false

    run dolt status
    [ $status -eq 0 ]
    [[ "$output" =~ "main" ]] || false

    run dolt checkout copy
    [ $status -eq 0 ]

    run dolt sql -q "SELECT * FROM test WHERE pk > 3;"
    [[ "$output" =~ "4" ]] || false
}

@test "sql-branch: DOLT_BRANCH -c throws error on error cases" {
    run dolt status
    [[ "$output" =~ "main" ]] || false

    # branch copying from is empty
    run dolt sql -q "SELECT DOLT_BRANCH('-c','','copy')"
    [ $status -eq 1 ]
    [ "$output" = "error: cannot branch empty string" ]

    # branch copying to is empty
    run dolt sql -q "SELECT DOLT_BRANCH('-c','main','')"
    [ $status -eq 1 ]
    [ "$output" = "error: cannot branch empty string" ]

    dolt branch 'existing_branch'
    run dolt branch
    [[ "$output" =~ "main" ]] || false
    [[ "$output" =~ "existing_branch" ]] || false
    [[ ! "$output" =~ "original" ]] || false

    # branch copying from that don't exist
    run dolt sql -q "SELECT DOLT_BRANCH('-c', 'original', 'copy');"
    [ $status -eq 1 ]
    [ "$output" = "fatal: A branch named 'original' not found" ]

    # branch copying to that exists
    run dolt sql -q "SELECT DOLT_BRANCH('-c', 'main', 'existing_branch');"
    [ $status -eq 1 ]
    [ "$output" = "fatal: A branch named 'existing_branch' already exists." ]
}

@test "sql-branch: DOLT_BRANCH works as insert into dolt_branches table" {
    dolt add . && dolt commit -m "1, 2, and 3 in test table"

    run dolt sql -q "SELECT hash FROM dolt_branches WHERE name='main';"
    [ $status -eq 0 ]
    mainhash=$output

    dolt sql -q "SELECT DOLT_BRANCH('feature-branch');"
    run dolt sql -q "SELECT hash FROM dolt_branches WHERE name='feature-branch';"
    [ $status -eq 0 ]
    [ "$output" = "$mainhash" ]
}

@test "sql-branch: DOLT_BRANCH -m renames current branch" {
    skip "need to handle renaming checked out branch on sql"
    dolt add . && dolt commit -m "0, 1, and 2 in test table"
    run dolt status
    [[ "$output" =~ "main" ]] || false

    # Current branch should be still main with test table without entry 4
    dolt sql << SQL
SELECT DOLT_BRANCH('-m', 'main', 'renamed');
INSERT INTO test VALUES (3);
SELECT DOLT_COMMIT('-am','add 3');
SELECT count(*) FROM dolt_branches;
SQL
    [ $status -eq 0 ]
    [[ ! "$output" =~ "2" ]] || false

    run dolt status
    [ $status -eq 0 ]
    [[ "$output" =~ "renamed" ]] || false
    [[ ! "$output" =~ "main" ]] || false
}

@test "sql-branch: DOLT_BRANCH -m renames branch not checked out" {
    dolt add . && dolt commit -m "0, 1, and 2 in test table"
    run dolt status
    [[ "$output" =~ "main" ]] || false

    dolt branch 'original'

    # Current branch should be still main with test table without entry 4
    run dolt sql << SQL
SELECT DOLT_BRANCH('-m', 'original', 'renamed');
SELECT count(*) FROM dolt_branches;
SQL
    [ $status -eq 0 ]
    [[ "$output" =~ "2" ]] || false

    run dolt branch
    [ $status -eq 0 ]
    [[ "$output" =~ "renamed" ]] || false
    [[ ! "$output" =~ "original" ]] || false
}

@test "sql-branch: DOLT_BRANCH -d works deleting a branch" {
    dolt add . && dolt commit -m "0, 1, and 2 in test table"
    run dolt status
    [[ "$output" =~ "main" ]] || false

    dolt branch new_branch

    # Current branch should be still main with test table without entry 4
    run dolt sql << SQL
SELECT DOLT_BRANCH('-d','new_branch');
SELECT count(*) FROM dolt_branches;
SQL
    [ $status -eq 0 ]
    [[ ! "$output" =~ "2" ]] || false

    run dolt branch
    [ $status -eq 0 ]
    [[ ! "$output" =~ "new_branch" ]] || false
}

@test "sql-branch: DOLT_BRANCH -d works deleting multiple branches" {
    dolt add . && dolt commit -m "0, 1, and 2 in test table"
    run dolt status
    [[ "$output" =~ "main" ]] || false

    dolt branch branch_one
    dolt branch branch_two

    # Current branch should be still main with test table without entry 4
    run dolt sql << SQL
SELECT DOLT_BRANCH('-d','branch_one','branch_two');
SELECT count(*) FROM dolt_branches;
SQL
    [ $status -eq 0 ]
    [[ "$output" =~ "1" ]] || false

    run dolt branch
    [ $status -eq 0 ]
    [[ ! "$output" =~ "branch_one" ]] || false
    [[ ! "$output" =~ "branch_two" ]] || false
}
