#!/usr/bin/env bats

load helpers

function setup() {
	teardown_busybox
	setup_busybox

	# Create fake rootfs.
	mkdir rootfs/testdir
	echo "Forbidden information!" > rootfs/testfile

	# add extra masked paths
	sed -i 's;"maskedPaths": \[;"maskedPaths": \["/testdir","/testfile",;g' config.json
}

function teardown() {
	teardown_busybox
}

@test "mask paths [file]" {
  skip "runv doesn't support mask paths yet"
	# run busybox detached
	runv run -d --console-socket $CONSOLE_SOCKET test_busybox
	[ "$status" -eq 0 ]

	runv exec test_busybox cat /testfile
	[ "$status" -eq 0 ]
	[[ "${output}" == "" ]]

	runv exec test_busybox rm -f /testfile
	[ "$status" -eq 1 ]
	[[ "${output}" == *"Read-only file system"* ]]

	runv exec test_busybox umount /testfile
	[ "$status" -eq 1 ]
	[[ "${output}" == *"Operation not permitted"* ]]
}

@test "mask paths [directory]" {
  skip "runv doesn't support mask paths yet"
	# run busybox detached
	runv run -d --console-socket $CONSOLE_SOCKET test_busybox
	[ "$status" -eq 0 ]

	runv exec test_busybox ls /testdir
	[ "$status" -eq 0 ]
	[[ "${output}" == "" ]]

	runv exec test_busybox touch /testdir/foo
	[ "$status" -eq 1 ]
	[[ "${output}" == *"Read-only file system"* ]]

	runv exec test_busybox rm -rf /testdir
	[ "$status" -eq 1 ]
	[[ "${output}" == *"Read-only file system"* ]]

	runv exec test_busybox umount /testdir
	[ "$status" -eq 1 ]
	[[ "${output}" == *"Operation not permitted"* ]]
}
