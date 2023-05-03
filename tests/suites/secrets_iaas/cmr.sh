run_secrets_cmr() {
	echo

	# First set up a cross model relation.
	juju --show-log add-model "model-secrets-offer"
	juju --show-log deploy juju-qa-dummy-source
	juju --show-log offer dummy-source:sink
	wait_for "dummy-source" "$(idle_condition "dummy-source")"

	juju --show-log add-model "model-secrets-consume"
	juju --show-log deploy juju-qa-dummy-sink
	juju --show-log integrate dummy-sink model-secrets-offer.dummy-source

	wait_for "dummy-sink" "$(idle_condition "dummy-sink")"
	wait_for "dummy-source" '.applications["dummy-sink"] | .relations.source[0]'

	juju switch "model-secrets-offer"
	juju config dummy-source token=foobar
	juju switch "model-secrets-consume"
	wait_for "active" '."application-endpoints"["dummy-source"]."application-status".current'

	juju switch "model-secrets-offer"
	wait_for "1" '.offers["dummy-source"]["active-connected-count"]'

	# Create and share a secret on the offer side.
	secret_uri=$(juju exec --unit dummy-source/0 -- secret-add foo=bar | cut -d: -f2)
	relation_id=$(juju --show-log show-unit -m model-secrets-offer dummy-source/0 --format json | jq '."dummy-source/0"."relation-info"[0]."relation-id"')
	juju exec --unit dummy-source/0 -- secret-grant "$secret_uri" -r "$relation_id"

	# Check the secret can be read by the consumer.
	juju switch "model-secrets-consume"
	# secret-get by URI - consume content.
	juju exec --unit dummy-sink/0 -- secret-get --label mylabel "$secret_uri" | grep 'foo: bar'
	# secret-get by URI - consume content.
	juju exec --unit dummy-sink/0 -- secret-get --label mylabel | grep 'foo: bar'

	# Add a new revision and check consumer can see it.
	juju switch "model-secrets-offer"
	juju exec --unit dummy-source/0 -- secret-set "$secret_uri" foo=bar2
	juju switch "model-secrets-consume"
	juju exec --unit dummy-sink/0 -- secret-get --label mylabel | grep 'foo: bar'
	juju exec --unit dummy-sink/0 -- secret-get --label mylabel --peek | grep 'foo: bar2'
	juju exec --unit dummy-sink/0 -- secret-get --label mylabel | grep 'foo: bar'
	juju exec --unit dummy-sink/0 -- secret-get --label mylabel --refresh | grep 'foo: bar2'
	juju exec --unit dummy-sink/0 -- secret-get --label mylabel | grep 'foo: bar2'

	# Suspend relation and check access is lost.
	juju switch "model-secrets-offer"
	juju suspend-relation "$relation_id"
	juju switch "model-secrets-consume"
	check_contains "$(juju exec --unit dummy-sink/0 -- secret-get "$secret_uri" 2>&1)" 'permission denied'
	# Resume relation and access is restored.
	juju switch "model-secrets-offer"
	juju resume-relation "$relation_id"
	juju switch "model-secrets-consume"
	juju exec --unit dummy-sink/0 -- secret-get --label mylabel | grep 'foo: bar2'

	# Secret-revoke by relation ID.
	juju switch "model-secrets-offer"
	juju exec --unit dummy-source/0 -- secret-revoke "$secret_uri" --relation "$relation_id"
	juju switch "model-secrets-consume"
	check_contains "$(juju exec --unit dummy-sink/0 -- secret-get "$secret_uri" 2>&1)" 'permission denied'
}

test_secrets_cmr() {
	if [ "$(skip 'test_secrets_cmr')" ]; then
		echo "==> TEST SKIPPED: test_secrets_cmr"
		return
	fi

	(
		set_verbosity

		cd .. || exit

		run "run_secrets_cmr"
	)
}
