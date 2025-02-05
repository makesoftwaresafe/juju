"""Check the Go package doc comments."""

import multiprocessing.pool
import os
import subprocess
import sys


def check_all():
	"""Check the package doc comment for all packages."""
	packages = find_packages('.')

	# Use a thread pool with os.cpu_count() threads to parallelize the
	# "go doc" invocations. Brings runtime down from 1m20s to 30s or so.
	pool = multiprocessing.pool.ThreadPool()
	errors = pool.map(check_package, packages)

	exit_error = False
	for package, error in zip(packages, errors):
		if error:
			print(package+': '+error)
			exit_error = True

	if exit_error:
		sys.exit(1)


def find_packages(path):
	"""Return a list of all Go packages under the given path."""
	packages = []
	for root, _, files in os.walk(path):
		has_go = any(f.endswith('.go') for f in files)
		if not has_go:
			continue
		packages.append(root)
	packages.sort()
	return packages


def check_package(package):
	"""Check the package doc comment for the given package. Ensure that if the
	package has a package doc comment, it's only in doc.go.
	"""
	if package.endswith('/mocks'):
		return None

	process = subprocess.run(['go', 'doc', package], capture_output=True, check=True)
	output = process.stdout.decode('utf-8').strip()

	# Reduce output to only the package doc comment (up to declarations)
	lines = []
	for line in output.splitlines():
		if line.startswith(('const ', 'func ', 'type ', 'var ')):
			break
		lines.append(line)
	header = '\n'.join(lines).strip()

	if '\n' not in header:
		# No package doc comment
		return None

	errors = []

	comment_files = package_comment_files(package)
	non_doc_files = [f for f in comment_files if f != 'doc.go']
	if non_doc_files:
		errors.append('package comment in non-doc.go files: '+', '.join(non_doc_files))

	lines = header.splitlines()
	if len(lines) > 2 and not lines[2].startswith('Package '):
		errors.append('package comment does not start with "Package "')

	if 'Copyright 20' in header:
		errors.append('package comment should not include Copyright notice (add blank line after)')

	return '; '.join(errors)


def package_comment_files(package):
	"""Return a list of the files in package that have package doc comments."""
	files = []
	for name in os.listdir(package):
		if not name.endswith('.go') or name.endswith('_test.go'):
			continue
		lines = []
		with open(os.path.join(package, name)) as f:
			for line in f:
				if line.startswith('package '):
					break
				lines.append(line)
		if lines and lines[-1].startswith('//'):
			files.append(name)
	return files


if __name__ == '__main__':
	check_all()
