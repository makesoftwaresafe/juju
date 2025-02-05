test_constraints_common() {
	if [ "$(skip 'test_constraints_common')" ]; then
		echo "==> TEST SKIPPED: constraints"
		return
	fi

	(
		set_verbosity

		cd .. || exit

		case "${BOOTSTRAP_PROVIDER:-}" in
		"lxd" | "lxd-remote" | "localhost")
			run "run_constraints_lxd"
			;;
		"microk8s")
			echo "==> TEST SKIPPED: constraints - there are no test for k8s cloud"
			;;
		*)
			run "run_constraints_vm"
			;;
		esac
	)
}
